package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/usecase"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 54 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Is this a device client or viewer client
	IsDevice bool
	DeviceID string

	// Usecase for processing incoming telemetry (only used by device clients)
	TelemetryUsecase usecase.TelemetryUsecase
}

// readPump pumps messages from the websocket connection to the hub/usecase.
func (c *Client) ReadPump() {
	defer func() {
		if c.IsDevice {
			c.hub.UnregisterDevice <- c
		} else {
			c.hub.UnregisterViewer <- c
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		if c.IsDevice && c.TelemetryUsecase != nil {
			var payload domain.TelemetryPayload
			if err := json.Unmarshal(message, &payload); err != nil {
				log.Printf("error unmarshaling payload: %v", err)
				continue
			}

			// Force device ID from connection authentication to avoid spoofing
			payload.DeviceID = c.DeviceID

			// Process the payload through usecase
			// The usecase will handle broadcasting to the hub via the injected channel
			if err := c.TelemetryUsecase.ProcessIncoming(context.Background(), payload); err != nil {
				log.Printf("error processing payload: %v", err)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

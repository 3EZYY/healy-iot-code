package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rafif/healy-backend/internal/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, check origin appropriately
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeDeviceWs handles websocket requests from the ESP32 devices.
func ServeDeviceWs(hub *Hub, telemetryUsecase usecase.TelemetryUsecase, w http.ResponseWriter, r *http.Request) {
	// The device_id should be extracted from the header as per blueprint
	deviceID := r.Header.Get("device_id")
	if deviceID == "" {
		// Fallback to query param if needed, or reject
		deviceID = r.URL.Query().Get("device_id")
		if deviceID == "" {
			http.Error(w, "Missing device_id", http.StatusBadRequest)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		hub:              hub,
		conn:             conn,
		send:             make(chan []byte, 256),
		IsDevice:         true,
		DeviceID:         deviceID,
		TelemetryUsecase: telemetryUsecase,
	}

	client.hub.RegisterDevice <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

// ServeViewerWs handles websocket requests from the frontend viewers.
func ServeViewerWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Assume JWT validation happened in middleware before reaching this handler
	// Or validate JWT from query param if required

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		IsDevice: false,
	}

	client.hub.RegisterViewer <- client

	go client.WritePump()
	go client.ReadPump()
}

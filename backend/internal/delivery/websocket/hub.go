package websocket

import "log"

// Hub maintains the set of active clients and broadcasts messages to the viewer clients.
type Hub struct {
	// Registered viewer clients (e.g., Frontend).
	viewerClients map[*Client]bool

	// Registered device clients (e.g., ESP32).
	deviceClients map[string]*Client // key is device_id

	// Inbound messages to broadcast to viewer clients.
	Broadcast chan []byte

	// Register requests from the clients.
	RegisterViewer chan *Client
	RegisterDevice chan *Client

	// Unregister requests from clients.
	UnregisterViewer chan *Client
	UnregisterDevice chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:        make(chan []byte, 256),
		RegisterViewer:   make(chan *Client),
		RegisterDevice:   make(chan *Client),
		UnregisterViewer: make(chan *Client),
		UnregisterDevice: make(chan *Client),
		viewerClients:    make(map[*Client]bool),
		deviceClients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterViewer:
			h.viewerClients[client] = true
			log.Println("Viewer client registered")

		case client := <-h.RegisterDevice:
			h.deviceClients[client.DeviceID] = client
			log.Printf("Device client registered: %s\n", client.DeviceID)

		case client := <-h.UnregisterViewer:
			if _, ok := h.viewerClients[client]; ok {
				delete(h.viewerClients, client)
				close(client.send)
				log.Println("Viewer client unregistered")
			}

		case client := <-h.UnregisterDevice:
			if c, ok := h.deviceClients[client.DeviceID]; ok && c == client {
				delete(h.deviceClients, client.DeviceID)
				close(client.send)
				log.Printf("Device client unregistered: %s\n", client.DeviceID)
			}

		case message := <-h.Broadcast:
			// Broadcast message to all viewer clients
			for client := range h.viewerClients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.viewerClients, client)
				}
			}
		}
	}
}

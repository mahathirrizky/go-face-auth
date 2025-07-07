package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader is a shared WebSocket upgrader for all handlers.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now. In production, you should restrict this.
		return true
	},
}

// Client represents a single WebSocket client connection.
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	CompanyID int // To identify which company this client belongs to
	Done chan struct{} // Channel to signal when the client is done
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	Register chan *Client
	Unregister chan *Client
	mu sync.RWMutex // Mutex to protect clients map
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %v (Company ID: %d)", client.Conn.RemoteAddr(), client.CompanyID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %v (Company ID: %d)", client.Conn.RemoteAddr(), client.CompanyID)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to client %v: %v", c.Conn.RemoteAddr(), err)
			return
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub.
// For dashboard updates, clients typically don't send messages, but this is good practice.
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
		close(c.Done) // Signal that this client is done
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from client %v: %v", c.Conn.RemoteAddr(), err)
			}
			break
		}
		// If clients send messages, process them here. For dashboard, usually not needed.
	}
}

// SendDashboardUpdate sends a dashboard summary update to all clients of a specific company.
func (h *Hub) SendDashboardUpdate(companyID int, summary interface{}) {
	message, err := json.Marshal(summary)
	if err != nil {
		log.Printf("Error marshalling dashboard summary: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.CompanyID == companyID {
			select {
			case client.Send <- message:
			default:
				// Client's Send buffer is full, likely a slow consumer. Close connection.
				close(client.Send)
				delete(h.clients, client)
				log.Printf("Closing slow client connection: %v (Company ID: %d)", client.Conn.RemoteAddr(), client.CompanyID)
			}
		}
	}
}

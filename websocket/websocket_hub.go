package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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

// CompanyBroadcastMessage represents a message to be broadcast to clients of a specific company.
type CompanyBroadcastMessage struct {
	CompanyID int
	Message   []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte // For general broadcasts (e.g., superadmin dashboard)
	Register chan *Client
	Unregister chan *Client
	BroadcastToCompany chan CompanyBroadcastMessage // New channel for company-specific broadcasts
	mu sync.RWMutex // Mutex to protect clients map
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		BroadcastToCompany: make(chan CompanyBroadcastMessage),
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
		case msg := <-h.BroadcastToCompany:
			h.mu.RLock()
			for client := range h.clients {
				if client.CompanyID == msg.CompanyID {
					select {
					case client.Send <- msg.Message:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessageToCompany sends a structured message to all clients of a specific company.
func (h *Hub) BroadcastMessageToCompany(companyID int, messageType string, payload interface{}) {
	structuredMessage := map[string]interface{}{
		"type":    messageType,
		"payload": payload,
	}
	messageBytes, err := json.Marshal(structuredMessage)
	if err != nil {
		log.Printf("Error marshalling broadcast message: %v", err)
		return
	}

	h.BroadcastToCompany <- CompanyBroadcastMessage{
		CompanyID: companyID,
		Message:   messageBytes,
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub.
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
		close(c.Done) // Signal that this client is done
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from client %v: %v", c.Conn.RemoteAddr(), err)
			}
			break
		}
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

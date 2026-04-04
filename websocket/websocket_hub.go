package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

var (
	allowedWebsocketOrigins []string
)

func init() {
	allowedWebsocketOrigins = parseAllowedOrigins(os.Getenv("WS_ALLOWED_ORIGINS"))
	if len(allowedWebsocketOrigins) == 0 {
		log.Println("WS_ALLOWED_ORIGINS not set, allowing all websocket origins")
	}
}

// Upgrader is a shared WebSocket upgrader for all handlers.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if len(allowedWebsocketOrigins) == 0 {
			return true
		}
		origin := r.Header.Get("Origin")
		return isOriginAllowed(origin)
	},
}

func parseAllowedOrigins(value string) []string {
	var origins []string
	if value == "" {
		return origins
	}
	for _, item := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func isOriginAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	for _, allowed := range allowedWebsocketOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

// Client represents a single WebSocket client connection.
type Client struct {
	Conn      *websocket.Conn
	Send      chan []byte
	CompanyID int           // To identify which company this client belongs to
	Done      chan struct{} // Channel to signal when the client is done
}

// CompanyBroadcastMessage represents a message to be broadcast to clients of a specific company.
type CompanyBroadcastMessage struct {
	CompanyID int
	Message   []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients             map[*Client]bool
	broadcast           chan []byte // For general broadcasts (e.g., superadmin dashboard)
	Register            chan *Client
	Unregister          chan *Client
	BroadcastToCompany  chan CompanyBroadcastMessage // New channel for company-specific broadcasts
	TriggerSuperAdminUpdate chan struct{}              // Channel to trigger super admin dashboard updates
	mu                  sync.RWMutex                 // Mutex to protect clients map
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:               make(chan []byte),
		Register:                make(chan *Client),
		Unregister:              make(chan *Client),
		BroadcastToCompany:      make(chan CompanyBroadcastMessage),
		TriggerSuperAdminUpdate: make(chan struct{}, 1), // Buffered so we don't block callers
		clients:                 make(map[*Client]bool),
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
			// Send message to all clients. If a client's send buffer is full,
			// mark it for removal and perform deletion under write lock.
			var toRemove []*Client
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
					// message sent
				default:
					// cannot send now, schedule removal
					toRemove = append(toRemove, client)
				}
			}
			h.mu.RUnlock()
			if len(toRemove) > 0 {
				h.mu.Lock()
				for _, client := range toRemove {
					if _, ok := h.clients[client]; ok {
						close(client.Send)
						delete(h.clients, client)
						log.Printf("Removed stuck client: %v (Company ID: %d)", client.Conn.RemoteAddr(), client.CompanyID)
					}
				}
				h.mu.Unlock()
			}
		case msg := <-h.BroadcastToCompany:
			var toRemove []*Client
			h.mu.RLock()
			for client := range h.clients {
				if client.CompanyID == msg.CompanyID {
					select {
					case client.Send <- msg.Message:
						// sent
					default:
						toRemove = append(toRemove, client)
					}
				}
			}
			h.mu.RUnlock()
			if len(toRemove) > 0 {
				h.mu.Lock()
				for _, client := range toRemove {
					if _, ok := h.clients[client]; ok {
						close(client.Send)
						delete(h.clients, client)
						log.Printf("Removed stuck client for company %d: %v", client.CompanyID, client.Conn.RemoteAddr())
					}
				}
				h.mu.Unlock()
			}
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

// SendSuperAdminDashboardUpdate broadcasts a pre-calculated dashboard summary to all superadmin clients.
func (h *Hub) SendSuperAdminDashboardUpdate(summary map[string]interface{}) {
	log.Println("Broadcasting superadmin dashboard update...")

	response := gin.H{
		"type":    "superadmin_dashboard_update",
		"payload": summary,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling superadmin dashboard update: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.CompanyID == 0 { // Target only superadmin clients
			select {
			case client.Send <- jsonResponse:
			default:
				log.Printf("Superadmin client send channel full or closed, removing client: %v", client.Conn.RemoteAddr())
			}
		}
	}
}

type SuperAdminNotificationPayload struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	CompanyID   uint   `json:"company_id,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
}

// SendSuperAdminNotification sends a structured notification to all superadmin clients.
func (h *Hub) SendSuperAdminNotification(payload SuperAdminNotificationPayload) {
	structuredMessage := map[string]interface{}{
		"type":    "superadmin_notification",
		"payload": payload,
	}
	messageBytes, err := json.Marshal(structuredMessage)
	if err != nil {
		log.Printf("Error marshalling superadmin notification: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.CompanyID == 0 { // Superadmin clients have CompanyID 0
			select {
			case client.Send <- messageBytes:
			default:
				log.Printf("Superadmin client send channel full or closed, removing client: %v", client.Conn.RemoteAddr())
				// In a real application, you might want to handle this more gracefully
				// For now, we just log and let the read/write pump handle disconnection
			}
		}
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

package handlers

import (
	"log"
	"net/http"

	"go-face-auth/middleware"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// EmployeeWebSocketHandler handles WebSocket connections for employees.
func EmployeeWebSocketHandler(hub *websocket.Hub, c *gin.Context) {
	// Extract token from query parameter
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Employee WebSocket: Missing token.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		log.Printf("WebSocket connection attempt with invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authentication token."})
		return
	}

	// Extract employeeID and companyID from claims
	employeeIDFloat, ok := claims["id"].(float64)
	if !ok {
		log.Println("Employee WebSocket: Employee ID not found in token claims.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Employee ID not found in token claims"})
		return
	}
	employeeID := int(employeeIDFloat)

	companyIDFloat, ok := claims["companyID"].(float64)
	if !ok {
		log.Println("Employee WebSocket: Company ID not found in token claims.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in token claims"})
		return
	}
	companyID := int(companyIDFloat)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Employee WebSocket upgrade failed:", err)
		return
	}

	// Create a new client and register it with the hub
	client := &websocket.Client{Conn: conn, Send: make(chan []byte, 256), CompanyID: companyID, Done: make(chan struct{})}
	hub.Register <- client

	log.Printf("Employee %d (Company ID: %d) WebSocket connected.", employeeID, companyID)

	// Start goroutines to handle reading and writing WebSocket messages
	go client.WritePump()
	go client.ReadPump(hub)
}

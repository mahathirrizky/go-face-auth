package handlers

import (
	"log"
	"net/http"

	"go-face-auth/middleware"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"

)

// ServeWs handles WebSocket requests for dashboard updates.
func ServeWs(hub *websocket.Hub, c *gin.Context) {
	// Extract token from query parameter
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("WebSocket connection attempt without token.")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication token missing."})
		return
	}

	// Validate token and get claims
	claims, err := middleware.ValidateToken(tokenString)
	if err != nil {
		log.Printf("WebSocket connection attempt with invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authentication token."})
		return
	}

	// Extract companyID from claims
	companyID, ok := claims["companyID"].(float64)
	if !ok {
		log.Println("Company ID not found or invalid in token claims for WebSocket connection.")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Company ID not found or invalid in token claims."})
		return
	}
	compID := int(companyID)

	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	client := &websocket.Client{Conn: conn, Send: make(chan []byte, 256), CompanyID: compID}
	hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump(hub)

	log.Printf("Dashboard WebSocket client connected for Company ID: %d", compID)
}

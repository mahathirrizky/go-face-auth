package handlers

import (
	"log"
	"net/http"

	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"

)

// ServeWs handles WebSocket requests for dashboard updates.
func ServeWs(hub *websocket.Hub, c *gin.Context) {
	// Get companyID from JWT claims set by AuthMiddleware
	companyID, exists := c.Get("companyID")
	if !exists {
		log.Println("Company ID not found in token claims for WebSocket connection.")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Company ID not found in token claims."})
		return
	}

	compIDFloat, ok := companyID.(float64)
	if !ok {
		log.Println("Invalid company ID type in token claims for WebSocket connection.")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid company ID type in token claims."})
		return
	}
	compID := int(compIDFloat)

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

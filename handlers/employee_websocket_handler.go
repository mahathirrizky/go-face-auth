package handlers

import (
	"log"
	"net/http"
	

	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect: `token.Method.(*jwt.SigningMethodHMAC)`
		// For simplicity, using a hardcoded secret. In production, use os.Getenv("JWT_SECRET")
		return []byte("supersecretjwtkeythatshouldbechangedinproduction"), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Employee WebSocket: Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Employee WebSocket: Invalid token claims.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
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

package handlers

import (
	"log"
	"net/http"

	"go-face-auth/helper"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// BroadcastMessageRequest represents the request body for broadcasting a message.
type BroadcastMessageRequest struct {
	Message    string `json:"message" binding:"required"`
	ExpireDate string `json:"expire_date"` // YYYY-MM-DD format
}

// BroadcastMessage handles broadcasting a message to all employees of a company.
func BroadcastMessage(hub *websocket.Hub, c *gin.Context) {
	var req BroadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	// Get CompanyID from JWT claims (assuming it's set by JWT middleware)
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}

	// Convert companyID to int (it comes as float64 from JWT claims)
	cidFloat, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid Company ID type. Expected float64.")
		return
	}
	cid := int(cidFloat)

	log.Printf("Admin of Company ID %d broadcasting message: %s (Expires: %s)", cid, req.Message, req.ExpireDate)

	// Broadcast the message using the WebSocket hub
	hub.BroadcastMessageToCompany(cid, "broadcast_message", gin.H{
		"message": req.Message,
		"expire_date": req.ExpireDate,
	})

	helper.SendSuccess(c, http.StatusOK, "Message broadcast successfully.", nil)
}


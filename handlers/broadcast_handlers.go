package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/services"

	"go-face-auth/websocket"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BroadcastHandler defines the interface for broadcast related handlers.
type BroadcastHandler interface {
	BroadcastMessage(hub *websocket.Hub, c *gin.Context)
	GetBroadcasts(c *gin.Context)
	MarkBroadcastAsRead(c *gin.Context)
}

// broadcastHandler is the concrete implementation of BroadcastHandler.
type broadcastHandler struct {
	broadcastService services.BroadcastService
}

// NewBroadcastHandler creates a new instance of BroadcastHandler.
func NewBroadcastHandler(broadcastService services.BroadcastService) BroadcastHandler {
	return &broadcastHandler{
		broadcastService: broadcastService,
	}
}

// BroadcastMessageRequest represents the request body for broadcasting a message.
type BroadcastMessageRequest struct {
	Message    string `json:"message" binding:"required"`
	ExpireDate string `json:"expire_date"` // YYYY-MM-DD format
}

// BroadcastMessage handles saving and broadcasting a message to all employees of a company.
func (h *broadcastHandler) BroadcastMessage(hub *websocket.Hub, c *gin.Context) {
	var req BroadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	companyIDFloat, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}
	companyID := uint(companyIDFloat.(float64))

	newMessage, err := h.broadcastService.CreateBroadcastMessage(companyID, req.Message, req.ExpireDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to save broadcast message: "+err.Error())
		return
	}

	log.Printf("Admin of Company ID %d broadcasting message: %s", companyID, req.Message)

	// Prepare payload for WebSocket broadcast
	// The `is_read` field will be false by default for a new message.
	payload := gin.H{
		"id":          newMessage.ID,
		"message":     newMessage.Message,
		"expire_date": newMessage.ExpireDate,
		"created_at":  newMessage.CreatedAt,
		"is_read":     false,
	}

	// Broadcast the message using the WebSocket hub
	hub.BroadcastMessageToCompany(int(companyID), "broadcast_message", payload)

	helper.SendSuccess(c, http.StatusOK, "Message broadcast successfully.", nil)
}

// GetBroadcasts retrieves all active broadcast messages for the logged-in employee.
func (h *broadcastHandler) GetBroadcasts(c *gin.Context) {
	companyIDFloat, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}
	companyID := uint(companyIDFloat.(float64))

	employeeIDFloat, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token.")
		return
	}
	employeeID := uint(employeeIDFloat.(float64))

	messages, err := h.broadcastService.GetBroadcastsForEmployee(companyID, employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve broadcast messages: "+err.Error())
		return
	}
	helper.SendSuccess(c, http.StatusOK, "Broadcast messages retrieved successfully.", messages)
}

// MarkBroadcastAsRead marks a specific broadcast message as read for the logged-in employee.
func (h *broadcastHandler) MarkBroadcastAsRead(c *gin.Context) {
	employeeIDFloat, exists := c.Get("id") // Assuming 'id' claim is employeeID for employees
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token.")
		return
	}
	employeeID := uint(employeeIDFloat.(float64))

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid message ID.")
		return
	}

	if err := h.broadcastService.MarkBroadcastAsRead(employeeID, uint(messageID)); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to mark message as read: "+err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Message marked as read.", nil)
}
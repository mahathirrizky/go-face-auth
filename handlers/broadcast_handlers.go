package handlers

import (
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BroadcastMessageRequest represents the request body for broadcasting a message.
type BroadcastMessageRequest struct {
	Message    string `json:"message" binding:"required"`
	ExpireDate string `json:"expire_date"` // YYYY-MM-DD format
}

// BroadcastMessage handles saving and broadcasting a message to all employees of a company.
func BroadcastMessage(hub *websocket.Hub, c *gin.Context) {
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

	var expireTime *time.Time
	if req.ExpireDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.ExpireDate)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid expire_date format. Use YYYY-MM-DD.")
			return
		}
		expireTime = &parsedTime
	}

	// Create and save the message to the database
	newMessage := models.BroadcastMessage{
		CompanyID:  companyID,
		Message:    req.Message,
		ExpireDate: expireTime,
	}

	broadcastRepo := repository.NewBroadcastRepository(database.DB)
	if err := broadcastRepo.Create(&newMessage); err != nil {
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

// GetBroadcasts retrieves all active broadcast messages for the logged-in employee
// and marks them as read upon retrieval.
func GetBroadcasts(c *gin.Context) {
	companyIDFloat, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}
	companyID := uint(companyIDFloat.(float64))

	employeeIDFloat, exists := c.Get("id") // Assuming 'id' claim is employeeID for employees
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token.")
		return
	}
	employeeID := uint(employeeIDFloat.(float64))

	broadcastRepo := repository.NewBroadcastRepository(database.DB)
	messages, err := broadcastRepo.GetForEmployee(companyID, employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve broadcast messages: "+err.Error())
		return
	}

	// Mark all retrieved messages as read for this employee
	for i := range messages {
		if !messages[i].IsRead {
			err := broadcastRepo.MarkAsRead(employeeID, messages[i].ID)
			if err != nil {
				log.Printf("Error marking message %d as read for employee %d: %v", messages[i].ID, employeeID, err)
				// Continue processing, but log the error
			}
			messages[i].IsRead = true // Optimistically update the in-memory status
		}
	}

	helper.SendSuccess(c, http.StatusOK, "Broadcast messages retrieved successfully.", messages)
}


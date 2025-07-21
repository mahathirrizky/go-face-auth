package handlers

import (
	"go-face-auth/database/repository"
	"go-face-auth/helper"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// InitialPasswordSetupRequest defines the structure for an initial password setup request.
type InitialPasswordSetupRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// InitialPasswordSetup handles the request to set the initial password for an employee.
func InitialPasswordSetup(c *gin.Context) {
	var req InitialPasswordSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate password complexity
	if !helper.IsValidPassword(req.Password) {
		helper.SendError(c, http.StatusBadRequest, "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number.")
		return
	}

	// Get token from database
	token, err := repository.GetPasswordResetToken(req.Token)
	if err != nil {
		log.Printf("Error retrieving initial password token %s: %v", req.Token, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to set initial password.")
		return
	}
	// Ensure it's an initial password token and not used/expired
	if token == nil || token.Used || token.ExpiresAt.Before(time.Now()) || token.TokenType != "employee_initial_password" {
		log.Printf("Invalid, expired, or used initial password token: %s (Used: %t, ExpiresAt: %s, Type: %s)", req.Token, token.Used, token.ExpiresAt.String(), token.TokenType)
		helper.SendError(c, http.StatusBadRequest, "Invalid or expired initial password setup token.")
		return
	}

	// Mark token as used immediately to prevent reuse
	if err := repository.MarkPasswordResetTokenAsUsed(token); err != nil {
		log.Printf("Error marking initial password token %s as used: %v", req.Token, err)
		// Continue, but log the error. The token is effectively used.
	}

	// Hash new password
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing new password for initial setup: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to hash password.")
		return
	}

	// Update employee password
	employee, err := repository.GetEmployeeByID(token.UserID)
	if err != nil || employee == nil {
		log.Printf("Initial password setup: Employee with ID %d not found or error: %v", token.UserID, err)
		helper.SendError(c, http.StatusNotFound, "User not found.")
		return
	}
	employee.Password = string(hashedPassword)
	if err := repository.UpdateEmployee(employee); err != nil {
		log.Printf("Error updating employee password for ID %d during initial setup: %v", token.UserID, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to update employee password.")
		return
	}
	// Set IsPasswordSet to true after successful initial password setup
	if err := repository.SetEmployeePasswordSet(uint(employee.ID), true); err != nil {
		log.Printf("Error setting IsPasswordSet for employee %d to true during initial setup: %v", employee.ID, err)
		// Log error but don't block response, as password is set
	}

	helper.SendSuccess(c, http.StatusOK, "Initial password has been set successfully.", nil)
}
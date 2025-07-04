package handlers

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ForgotPasswordRequest defines the structure for a forgot password request.
type ForgotPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
}

// ForgotPassword handles the request to initiate a password reset.
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate email format
	if !helper.IsValidEmail(req.Email) {
		log.Printf("Forgot password: Invalid email format for %s", req.Email)
		helper.SendError(c, http.StatusBadRequest, "Invalid email format.")
		return
	}

	var userID int
	var userEmail string

	// Find the admin user
	admin, err := repository.GetAdminCompanyByUsername(req.Email)
		if err != nil || admin == nil {
			log.Printf("Forgot password: Admin with email %s not found or error: %v", req.Email, err)
			helper.SendError(c, http.StatusNotFound, "Email not registered.") // Explicitly tell frontend
			return
		}
		userID = admin.ID
		userEmail = admin.Email

	// Generate a unique token
	tokenString := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour) // Token valid for 1 hour

	passwordResetToken := &models.PasswordResetTokenTable{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	// Save token to database
	if err := repository.CreatePasswordResetToken(passwordResetToken); err != nil {
		log.Printf("Error creating password reset token for %s: %v", req.Email, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate password reset link.")
		return
	}

	// Construct reset URL (Frontend URL)
	// This should be your frontend's reset password page, e.g., https://your-frontend.com/reset-password?token=YOUR_TOKEN
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", helper.GetFrontendBaseURL(), tokenString)

	// Send email with reset link in a goroutine
	go func() {
		if err := helper.SendPasswordResetEmail(userEmail, resetURL); err != nil {
			log.Printf("Error sending password reset email to %s: %v", userEmail, err)
		}
	}()

	helper.SendSuccess(c, http.StatusOK, "If an account with that email exists, a password reset link has been sent.", nil)
}

// ResetPasswordRequest defines the structure for a reset password request.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword handles the request to reset the password using a token.
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get token from database
	token, err := repository.GetPasswordResetToken(req.Token)
	if err != nil {
		log.Printf("Error retrieving password reset token %s: %v", req.Token, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to reset password.")
		return
	}
	if token == nil || token.Used || token.ExpiresAt.Before(time.Now()) {
		log.Printf("Invalid, expired, or used token: %s (Used: %t, ExpiresAt: %s)", req.Token, token.Used, token.ExpiresAt.String())
		helper.SendError(c, http.StatusBadRequest, "Invalid or expired password reset token.")
		return
	}

	// Mark token as used immediately to prevent reuse
	if err := repository.MarkPasswordResetTokenAsUsed(token); err != nil {
		log.Printf("Error marking token %s as used: %v", req.Token, err)
		// Continue, but log the error. The token is effectively used.
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing new password: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to hash password.")
		return
	}

	// Update admin's password
	admin, err := repository.GetAdminCompanyByID(token.UserID)
		if err != nil || admin == nil {
			log.Printf("Reset password: Admin with ID %d not found or error: %v", token.UserID, err)
			helper.SendError(c, http.StatusNotFound, "User not found.")
			return
		}
		admin.Password = string(hashedPassword)
		if err := repository.UpdateAdminCompany(admin); err != nil {
			log.Printf("Error updating admin password for ID %d: %v", token.UserID, err)
			helper.SendError(c, http.StatusInternalServerError, "Failed to update admin password.")
			return
		}

	helper.SendSuccess(c, http.StatusOK, "Password has been reset successfully.", nil)
}


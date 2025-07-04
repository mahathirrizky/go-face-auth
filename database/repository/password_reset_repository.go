package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CreatePasswordResetToken creates a new password reset token in the database.
func CreatePasswordResetToken(token *models.PasswordResetTokenTable) error {
	result := database.DB.Create(token)
	if result.Error != nil {
		log.Printf("Error creating password reset token: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetPasswordResetToken retrieves a password reset token by its token string.
func GetPasswordResetToken(tokenString string) (*models.PasswordResetTokenTable, error) {
	var token models.PasswordResetTokenTable
	result := database.DB.Where("token = ? AND expires_at > ? AND used = ?", tokenString, time.Now(), false).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Token not found or expired/used
		}
		log.Printf("Error getting password reset token: %v", result.Error)
		return nil, result.Error
	}
	return &token, nil
}

// MarkPasswordResetTokenAsUsed marks a password reset token as used.
func MarkPasswordResetTokenAsUsed(token *models.PasswordResetTokenTable) error {
	token.Used = true
	result := database.DB.Save(token)
	if result.Error != nil {
		log.Printf("Error marking password reset token as used: %v", result.Error)
		return result.Error
	}
	return nil
}

package repository

import (
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// CreatePasswordResetToken creates a new password reset token in the database.
func (r *passwordResetRepository) CreatePasswordResetToken(token *models.PasswordResetTokenTable) error {
	result := r.db.Create(token)
	if result.Error != nil {
		log.Printf("Error creating password reset token: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetPasswordResetToken retrieves a password reset token by its token string.
func (r *passwordResetRepository) GetPasswordResetToken(tokenString string) (*models.PasswordResetTokenTable, error) {
	var token models.PasswordResetTokenTable
	result := r.db.Where("token = ? AND expires_at > ? AND used = ?", tokenString, time.Now(), false).First(&token)
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
func (r *passwordResetRepository) MarkPasswordResetTokenAsUsed(token *models.PasswordResetTokenTable) error {
	token.Used = true
	result := r.db.Save(token)
	if result.Error != nil {
		log.Printf("Error marking password reset token as used: %v", result.Error)
		return result.Error
	}
	return nil
}

// InvalidatePasswordResetTokensByUserID marks all active password reset tokens for a user and type as used.
func (r *passwordResetRepository) InvalidatePasswordResetTokensByUserID(userID uint, tokenType string) error {
	result := r.db.Model(&models.PasswordResetTokenTable{}).
		Where("user_id = ? AND token_type = ? AND used = ?", userID, tokenType, false).
		Update("used", true)
	if result.Error != nil {
		log.Printf("Error invalidating password reset tokens for user %d, type %s: %v", userID, tokenType, result.Error)
		return result.Error
	}
	return nil
}
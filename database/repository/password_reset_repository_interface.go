package repository

import "go-face-auth/models"

// PasswordResetRepository defines the contract for password_reset-related database operations.
type PasswordResetRepository interface {
	CreatePasswordResetToken(token *models.PasswordResetTokenTable) error
	GetPasswordResetToken(tokenString string) (*models.PasswordResetTokenTable, error)
	MarkPasswordResetTokenAsUsed(token *models.PasswordResetTokenTable) error
	InvalidatePasswordResetTokensByUserID(userID uint, tokenType string) error
}

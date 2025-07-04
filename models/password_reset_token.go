package models

import (
	"time"

	"gorm.io/gorm"
)

// UserType defines the type of user for password reset tokens.
// PasswordResetTokenTable represents the password reset tokens in the database.
type PasswordResetTokenTable struct {
	gorm.Model
	Token      string    `gorm:"uniqueIndex;not null"`
	UserID     int       `gorm:"not null"` // ID of the admin user
	ExpiresAt  time.Time `gorm:"not null"`
	Used       bool      `gorm:"default:false"` // To prevent token reuse
}

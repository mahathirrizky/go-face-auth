package models

import (
	"time"

	"gorm.io/gorm"
)

// UserType defines the type of user for password reset tokens.
type UserType string

const (
	UserTypeEmployee UserType = "employee"
	UserTypeAdmin    UserType = "admin"
)

// PasswordResetTokenTable represents the password reset tokens in the database.
type PasswordResetTokenTable struct {
	gorm.Model
	Token      string    `gorm:"uniqueIndex;not null"`
		UserID     int       `gorm:"not null"` // ID of the user (employee or admin)
	UserType   UserType  `gorm:"type:varchar(20);not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	Used       bool      `gorm:"default:false"` // To prevent token reuse
}

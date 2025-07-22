package models

import (
	"time"

	"gorm.io/gorm"
)

// CustomPackageRequest represents a request for a custom subscription package.
type CustomPackageRequest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CompanyID   uint           `gorm:"not null" json:"company_id"`
	Email       string         `gorm:"not null" json:"email"`
	Phone       string         `json:"phone,omitempty"`
	CompanyName string         `gorm:"not null" json:"company_name"`
	Message     string         `gorm:"type:text" json:"message"`
	Status      string         `gorm:"default:'pending'" json:"status"` // e.g., pending, contacted, resolved
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
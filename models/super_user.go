package models

import (
	"time"

	"gorm.io/gorm"
)

// SuperUser represents the superuser model in the database.
type SuperUserTable struct {
	ID        int           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`	
	Password  string         `gorm:"not null" json:"-"` // Password should not be marshaled to JSON
	Role      string         `gorm:"not null" json:"role"` // e.g., "admin", "super_admin"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

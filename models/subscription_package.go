package models

import (
	"time"

	"gorm.io/gorm"
)

// SubscriptionPackageTable represents the subscription package model in the database.
type SubscriptionPackageTable struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"unique;not null" json:"name"`
	Price        float64        `gorm:"not null" json:"price"`
	DurationInMonths int `json:"duration_in_months"`
	MaxEmployees int            `gorm:"not null" json:"max_employees"`
	Features     string         `gorm:"type:text" json:"features"` // Comma-separated list of features or JSON string
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

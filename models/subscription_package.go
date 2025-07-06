package models

import (
	"time"

	"gorm.io/gorm"
)

// SubscriptionPackageTable represents the subscription package model in the database.
type SubscriptionPackageTable struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"unique;not null" json:"name"`
	PriceMonthly float64        `gorm:"not null" json:"price_monthly"`
	PriceYearly  float64        `gorm:"not null" json:"price_yearly"`
	MaxEmployees int            `gorm:"not null" json:"max_employees"`
	Features     string         `gorm:"type:text" json:"features"` // Comma-separated list of features or JSON string
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

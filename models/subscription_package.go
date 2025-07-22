package models

import (
	"time"

	"gorm.io/gorm"
)

// SubscriptionPackageTable represents the subscription package model in the database.
type SubscriptionPackageTable struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	PackageName  string         `gorm:"unique;not null" json:"package_name"`
	PriceMonthly float64        `gorm:"not null" json:"price_monthly"`
	PriceYearly  float64        `gorm:"not null" json:"price_yearly"`
	MaxEmployees int            `gorm:"not null" json:"max_employees"`
	MaxLocations int            `gorm:"not null;default:0" json:"max_locations"` // New field for max locations
	Features     string         `gorm:"type:text" json:"features"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
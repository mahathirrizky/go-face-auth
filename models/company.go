package models

import "time"

// Company represents a company in the system.
type CompaniesTable struct {
	ID                   int           `gorm:"primaryKey" json:"id"`
	Name                 string         `gorm:"not null" json:"name"`
	Address              string         `json:"address"`
	SubscriptionPackageID int           `json:"subscription_package_id"`
	SubscriptionPackage  SubscriptionPackageTable `gorm:"foreignKey:SubscriptionPackageID" json:"-"`
	SubscriptionStatus   string         `gorm:"not null;default:'pending'" json:"subscription_status"` // e.g., 'pending', 'active', 'trial', 'expired'
	SubscriptionStartDate *time.Time    `json:"subscription_start_date,omitempty"`
	SubscriptionEndDate   *time.Time    `json:"subscription_end_date,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	AdminCompaniesTable []AdminCompaniesTable `gorm:"foreignKey:CompanyID"` // Has many AdminCompaniesTable
}
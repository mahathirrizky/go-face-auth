package repository

import "go-face-auth/models"

// SubscriptionPackageRepository defines the contract for subscription_package-related database operations.
type SubscriptionPackageRepository interface {
	GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error)
	CreateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error
	UpdateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error
	UpdateSubscriptionPackageFields(id string, updates map[string]interface{}) error
	DeleteSubscriptionPackage(id string) error
	GetSubscriptionPackageByID(id int) (*models.SubscriptionPackageTable, error)
}

package services

import (
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

// SubscriptionPackageService defines the interface for subscription package business logic.
type SubscriptionPackageService interface {
	GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error)
	CreateSubscriptionPackage(newPackage *models.SubscriptionPackageTable) error
	UpdateSubscriptionPackage(packageID string, updates map[string]interface{}) error
	DeleteSubscriptionPackage(packageID string) error
}

// subscriptionPackageService is the concrete implementation of SubscriptionPackageService.
type subscriptionPackageService struct {
	subscriptionPackageRepo repository.SubscriptionPackageRepository
}

// NewSubscriptionPackageService creates a new instance of SubscriptionPackageService.
func NewSubscriptionPackageService(subscriptionPackageRepo repository.SubscriptionPackageRepository) SubscriptionPackageService {
	return &subscriptionPackageService{
		subscriptionPackageRepo: subscriptionPackageRepo,
	}
}

func (s *subscriptionPackageService) GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error) {
	return s.subscriptionPackageRepo.GetSubscriptionPackages()
}

func (s *subscriptionPackageService) CreateSubscriptionPackage(newPackage *models.SubscriptionPackageTable) error {
	return s.subscriptionPackageRepo.CreateSubscriptionPackage(newPackage)
}

func (s *subscriptionPackageService) UpdateSubscriptionPackage(packageID string, updates map[string]interface{}) error {
	return s.subscriptionPackageRepo.UpdateSubscriptionPackageFields(packageID, updates)
}

func (s *subscriptionPackageService) DeleteSubscriptionPackage(packageID string) error {
	return s.subscriptionPackageRepo.DeleteSubscriptionPackage(packageID)
}
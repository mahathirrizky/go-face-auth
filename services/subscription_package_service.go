package services

import (
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

func GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error) {
	return repository.GetSubscriptionPackages()
}

func CreateSubscriptionPackage(newPackage *models.SubscriptionPackageTable) error {
	return repository.CreateSubscriptionPackage(newPackage)
}

func UpdateSubscriptionPackage(packageID string, updates map[string]interface{}) error {
	return repository.UpdateSubscriptionPackageFields(packageID, updates)
}

func DeleteSubscriptionPackage(packageID string) error {
	return repository.DeleteSubscriptionPackage(packageID)
}

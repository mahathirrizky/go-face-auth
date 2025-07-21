package repository

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/models"
	"strconv"


)

// GetSubscriptionPackages retrieves all available subscription packages.
func GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error) {
	var packages []models.SubscriptionPackageTable
	if err := database.DB.Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

// CreateSubscriptionPackage creates a new subscription package.
func CreateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if err := database.DB.Create(pkg).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSubscriptionPackage updates an existing subscription package (full update).
func UpdateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if err := database.DB.Save(pkg).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSubscriptionPackageFields updates only the specified fields of a subscription package.
func UpdateSubscriptionPackageFields(id string, updates map[string]interface{}) error {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid package ID: %w", err)
	}

	if err := database.DB.Model(&models.SubscriptionPackageTable{}).Where("id = ?", uint(idUint)).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSubscriptionPackage deletes a subscription package.
func DeleteSubscriptionPackage(id string) error {
	if err := database.DB.Delete(&models.SubscriptionPackageTable{}, id).Error; err != nil {
		return err
	}
	return nil
}

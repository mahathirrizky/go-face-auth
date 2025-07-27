package repository

import (
	"fmt"

	"go-face-auth/models"
	"strconv"

	"gorm.io/gorm"
)

type subscriptionPackageRepository struct {
	db *gorm.DB
}

func NewSubscriptionPackageRepository(db *gorm.DB) SubscriptionPackageRepository {
	return &subscriptionPackageRepository{db: db}
}

// GetSubscriptionPackages retrieves all available subscription packages.
func (r *subscriptionPackageRepository) GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error) {
	var packages []models.SubscriptionPackageTable
	if err := r.db.Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

// CreateSubscriptionPackage creates a new subscription package.
func (r *subscriptionPackageRepository) CreateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if err := r.db.Create(pkg).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSubscriptionPackage updates an existing subscription package (full update).
func (r *subscriptionPackageRepository) UpdateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if err := r.db.Save(pkg).Error; err != nil {
		return err
	}
	return nil
}

// UpdateSubscriptionPackageFields updates only the specified fields of a subscription package.
func (r *subscriptionPackageRepository) UpdateSubscriptionPackageFields(id string, updates map[string]interface{}) error {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid package ID: %w", err)
	}

	if err := r.db.Model(&models.SubscriptionPackageTable{}).Where("id = ?", uint(idUint)).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSubscriptionPackage deletes a subscription package.
func (r *subscriptionPackageRepository) DeleteSubscriptionPackage(id string) error {
	if err := r.db.Delete(&models.SubscriptionPackageTable{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetSubscriptionPackageByID retrieves a subscription package by its ID.
func (r *subscriptionPackageRepository) GetSubscriptionPackageByID(id int) (*models.SubscriptionPackageTable, error) {
	var pkg models.SubscriptionPackageTable
	if err := r.db.First(&pkg, id).Error; err != nil {
		return nil, err
	}
	return &pkg, nil
}
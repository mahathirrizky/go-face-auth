package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"strings"

	"gorm.io/gorm"
)

// CreateCustomPackageRequest inserts a new custom package request into the database.
func CreateCustomPackageRequest(req *models.CustomPackageRequest) error {
	result := database.DB.Create(req)
	if result.Error != nil {
		log.Printf("Error creating custom package request: %v", result.Error)
		return result.Error
	}
	log.Printf("Custom package request created with ID: %d", req.ID)
	return nil
}

// GetCustomPackageRequestsPaginated retrieves custom package requests with pagination and search.
func GetCustomPackageRequestsPaginated(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error) {
	var requests []models.CustomPackageRequest
	var totalRecords int64

	query := database.DB.Model(&models.CustomPackageRequest{})

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(company_name) LIKE ? OR LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting custom package requests: %v", err)
		return nil, 0, err
	}

	// Apply pagination and order
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	if result.Error != nil {
		log.Printf("Error getting paginated custom package requests: %v", result.Error)
		return nil, 0, result.Error
	}

	return requests, totalRecords, nil
}

// GetCustomPackageRequestByID retrieves a custom package request by its ID.
func GetCustomPackageRequestByID(id uint) (*models.CustomPackageRequest, error) {
	var req models.CustomPackageRequest
	result := database.DB.First(&req, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		log.Printf("Error getting custom package request by ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &req, nil
}

// UpdateCustomPackageRequest updates an existing custom package request.
func UpdateCustomPackageRequest(req *models.CustomPackageRequest) error {
	result := database.DB.Save(req)
	if result.Error != nil {
		log.Printf("Error updating custom package request %d: %v", req.ID, result.Error)
		return result.Error
	}
	log.Printf("Custom package request %d updated successfully.", req.ID)
	return nil
}
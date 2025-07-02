package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateCompany inserts a new company into the database.
func CreateCompany(company *models.Company) error {
	result := database.DB.Create(company)
	if result.Error != nil {
		log.Printf("Error creating company: %v", result.Error)
		return result.Error
	}
	log.Printf("Company created with ID: %d", company.ID)
	return nil
}

// GetCompanyByID retrieves a company by its ID.
func GetCompanyByID(id int) (*models.Company, error) {
	var company models.Company
	result := database.DB.First(&company, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Company not found
		}
		log.Printf("Error getting company with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &company, nil
}

package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateCompany inserts a new company into the database.
func CreateCompany(company *models.CompaniesTable) error {
	result := database.DB.Create(company)
	if result.Error != nil {
		log.Printf("Error creating company: %v", result.Error)
		return result.Error
	}
	log.Printf("Company created with ID: %d", company.ID)
	return nil
}

// GetCompanyByID retrieves a company by its ID.
func GetCompanyByID(id int) (*models.CompaniesTable, error) {
	var company models.CompaniesTable
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

// UpdateCompany updates an existing company in the database.
func UpdateCompany(company *models.CompaniesTable) error {
	result := database.DB.Save(company)
	if result.Error != nil {
		log.Printf("Error updating company: %v", result.Error)
		return result.Error
	}
	log.Printf("Company with ID %d updated.", company.ID)
	return nil
}

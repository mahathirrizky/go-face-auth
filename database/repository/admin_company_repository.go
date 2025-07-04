package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateAdminCompany inserts a new AdminCompany record into the database.
func CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	result := database.DB.Create(adminCompany)
	if result.Error != nil {
		log.Printf("Error creating AdminCompany: %v", result.Error)
		return result.Error
	}
	log.Printf("AdminCompany created with ID: %d", adminCompany.ID)
	return nil
}

// GetAdminCompanyByCompanyID retrieves an AdminCompany record by CompanyID.
func GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := database.DB.Where("company_id = ?", companyID).First(&adminCompany)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found for this CompanyID
		}
		log.Printf("Error getting AdminCompany by CompanyID %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}

// GetAdminCompanyByEmployeeID retrieves an AdminCompany record by EmployeeID.
func GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := database.DB.Where("employee_id = ?", employeeID).First(&adminCompany)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found for this EmployeeID
		}
		log.Printf("Error getting AdminCompany by EmployeeID %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}

// GetAdminCompanyByUsername retrieves an AdminCompany record by Username.
func GetAdminCompanyByUsername(username string) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := database.DB.Preload("Company").Where("email = ?", username).First(&adminCompany)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found for this Username
		}
		log.Printf("Error getting AdminCompany by Username %s: %v", username, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}
package repository

import (
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

type adminCompanyRepository struct {
	db *gorm.DB
}

func NewAdminCompanyRepository(db *gorm.DB) AdminCompanyRepository {
	return &adminCompanyRepository{db: db}
}

// CreateAdminCompany inserts a new AdminCompany record into the database.
func (r *adminCompanyRepository) CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	result := r.db.Create(adminCompany)
	if result.Error != nil {
		log.Printf("Error creating AdminCompany: %v", result.Error)
		return result.Error
	}
	log.Printf("AdminCompany created with ID: %d", adminCompany.ID)
	return nil
}

// GetAdminCompanyByCompanyID retrieves an AdminCompany record by CompanyID.
func (r *adminCompanyRepository) GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := r.db.Where("company_id = ?", companyID).First(&adminCompany)
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
func (r *adminCompanyRepository) GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := r.db.Where("employee_id = ?", employeeID).First(&adminCompany)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found for this EmployeeID
		}
		log.Printf("Error getting AdminCompany by EmployeeID %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}

// GetAdminCompanyByEmail retrieves an AdminCompany record by Email.
func (r *adminCompanyRepository) GetAdminCompanyByEmail(email string) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := r.db.Preload("Company").Where("email = ?", email).First(&adminCompany)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found for this Email
		}
		log.Printf("Error getting AdminCompany by Email %s: %v", email, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}

// GetAdminCompanyByID retrieves an AdminCompany record by its ID.
func (r *adminCompanyRepository) GetAdminCompanyByID(id int) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	result := r.db.First(&adminCompany, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // AdminCompany not found
		}
		log.Printf("Error getting AdminCompany by ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &adminCompany, nil
}

// UpdateAdminCompany updates an existing AdminCompany record in the database.
func (r *adminCompanyRepository) UpdateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	result := r.db.Save(adminCompany)
	if result.Error != nil {
		log.Printf("Error updating AdminCompany: %v", result.Error)
		return result.Error
	}
	log.Printf("AdminCompany updated with ID: %d", adminCompany.ID)
	return nil
}

// ChangeAdminPassword updates the password for a specific admin company user.
func (r *adminCompanyRepository) ChangeAdminPassword(adminID int, newPasswordHash string) error {
	result := r.db.Model(&models.AdminCompaniesTable{}).Where("id = ?", adminID).Update("password", newPasswordHash)
	if result.Error != nil {
		log.Printf("Error updating admin password for admin ID %d: %v", adminID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("No admin found with ID %d to update password", adminID)
		return gorm.ErrRecordNotFound // Or a custom error
	}
	log.Printf("Password updated successfully for admin ID %d", adminID)
	return nil
}
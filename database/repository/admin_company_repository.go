package repository

import (
	"go-face-auth/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type adminCompanyRepository struct {
	db *gorm.DB
}

func NewAdminCompanyRepository(db *gorm.DB) AdminCompanyRepository {
	return &adminCompanyRepository{db: db}
}

// CreateAdminCompany inserperiots a new AdminCompany record into the database.
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
func (r *adminCompanyRepository) GetAdminCompanyByConfirmationToken(token string) (*models.AdminCompaniesTable, error) {
	var adminCompany models.AdminCompaniesTable
	if err := r.db.Where("confirmation_token = ?", token).First(&adminCompany).Error; err != nil {
		return nil, err
	}
	return &adminCompany, nil
}

func (r *adminCompanyRepository) UpdateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	return r.db.Save(adminCompany).Error
}

func (r *adminCompanyRepository) ChangeAdminPassword(adminID uint, newPassword string) error {
	admin, err := r.GetAdminCompanyByID(int(adminID))
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.Password = string(hashedPassword)
	return r.db.Save(admin).Error
}
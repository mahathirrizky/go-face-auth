package repository

import (
	"go-face-auth/models"
	"gorm.io/gorm"
)

// companyRepository is the concrete implementation of CompanyRepository
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository creates a new instance of CompanyRepository.
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) CreateCompany(company *models.CompaniesTable) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) CreateCompanyWithAdminAndShift(company *models.CompaniesTable, admin *models.AdminCompaniesTable, shift *models.ShiftsTable) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(company).Error; err != nil {
		tx.Rollback()
		return err
	}

	admin.CompanyID = company.ID
	if err := tx.Create(admin).Error; err != nil {
		tx.Rollback()
		return err
	}

	shift.CompanyID = company.ID
	if err := tx.Create(shift).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *companyRepository) GetAllCompanies() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	err := r.db.Find(&companies).Error
	return companies, err
}

func (r *companyRepository) DeleteCompany(id int) error {
	return r.db.Delete(&models.CompaniesTable{}, id).Error
}

func (r *companyRepository) GetCompanyByID(id int) (*models.CompaniesTable, error) {
	var company models.CompaniesTable
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanyWithSubscriptionDetails preloads subscription and offer details for a company.
func (r *companyRepository) GetCompanyWithSubscriptionDetails(id int) (*models.CompaniesTable, error) {
	var company models.CompaniesTable
	err := r.db.Preload("SubscriptionPackage").Preload("CustomOffer").First(&company, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err // Return error here as the caller expects a company
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) UpdateCompany(company *models.CompaniesTable) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) GetTotalEmployeesByCompanyID(companyID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EmployeesTable{}).Where("company_id = ?", companyID).Count(&count).Error
	return count, err
}

func (r *companyRepository) GetAllActiveCompanies() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	err := r.db.Where("subscription_status = ? OR subscription_status = ?", "active", "trial").Find(&companies).Error
	return companies, err
}

func (r *companyRepository) IsCompanyNameTaken(companyName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.CompaniesTable{}).Where("name = ?", companyName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
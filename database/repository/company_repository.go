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

func (r *companyRepository) GetCompanyByID(id int) (*models.CompaniesTable, error) {
	var company models.CompaniesTable
	if err := r.db.First(&company, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error for a single get
		}
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

func (r *companyRepository) GetAllActiveCompanies() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	err := r.db.Where("subscription_status = ? OR subscription_status = ?", "active", "trial").Find(&companies).Error
	return companies, err
}
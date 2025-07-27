package repository

import "go-face-auth/models"

// CompanyRepository defines the contract for company-related database operations.
type CompanyRepository interface {
	CreateCompany(company *models.CompaniesTable) error
	GetCompanyByID(id int) (*models.CompaniesTable, error)
	GetCompanyWithSubscriptionDetails(id int) (*models.CompaniesTable, error)
	UpdateCompany(company *models.CompaniesTable) error
	GetAllActiveCompanies() ([]models.CompaniesTable, error)
}

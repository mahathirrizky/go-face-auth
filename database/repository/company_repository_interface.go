package repository

import "go-face-auth/models"

// CompanyRepository defines the contract for company-related database operations.
type CompanyRepository interface {
	CreateCompany(company *models.CompaniesTable) error
	GetCompanyByID(id int) (*models.CompaniesTable, error)
	GetAllCompanies() ([]models.CompaniesTable, error)
	UpdateCompany(company *models.CompaniesTable) error
	DeleteCompany(id int) error
	GetCompanyWithSubscriptionDetails(companyID int) (*models.CompaniesTable, error)
	GetTotalEmployeesByCompanyID(companyID int) (int64, error)
	GetAllActiveCompanies() ([]models.CompaniesTable, error)
	CreateCompanyWithAdminAndShift(company *models.CompaniesTable, admin *models.AdminCompaniesTable, shift *models.ShiftsTable) error
	IsCompanyNameTaken(companyName string) (bool, error)
}

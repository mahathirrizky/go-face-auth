package repository

import "go-face-auth/models"

// AdminCompanyRepository defines the contract for admin_company-related database operations.
type AdminCompanyRepository interface {
	CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error
	GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByEmail(email string) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByID(id int) (*models.AdminCompaniesTable, error)
	UpdateAdminCompany(adminCompany *models.AdminCompaniesTable) error
	ChangeAdminPassword(adminID int, newPasswordHash string) error
}

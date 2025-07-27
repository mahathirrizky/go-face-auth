package repository

import "go-face-auth/models"

// SuperAdminRepository defines the contract for super_admin-related database operations.
type SuperAdminRepository interface {
	CreateSuperAdmin(superUser *models.SuperAdminTable) error
	GetSuperAdminByID(id int) (*models.SuperAdminTable, error)
	GetSuperAdminByEmail(email string) (*models.SuperAdminTable, error)
	UpdateSuperAdmin(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error)
	DeleteSuperAdmin(id int) error
	GetAllSuperAdmins() ([]models.SuperAdminTable, error)
}

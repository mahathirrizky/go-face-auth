package repository

import (
	"go-face-auth/models"
	"time"
)

// SuperAdminRepository defines the contract for super_admin-related database operations.
type SuperAdminRepository interface {
	CreateSuperAdmin(superUser *models.SuperAdminTable) error
	GetSuperAdminByID(id int) (*models.SuperAdminTable, error)
	GetSuperAdminByEmail(email string) (*models.SuperAdminTable, error)
	UpdateSuperAdmin(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error)
	DeleteSuperAdmin(id int) error
	GetAllSuperAdmins() ([]models.SuperAdminTable, error)

	// New methods for dashboard summary
	GetTotalCompaniesCount() (int64, error)
	GetCompaniesCountBySubscriptionStatus(status string) (int64, error)
	GetExpiredAndTrialExpiredCompaniesCount() (int64, error)
	GetRecentCompanies(limit int) ([]models.CompaniesTable, error)

	// New methods for companies and subscriptions
	GetAllCompaniesWithPreload() ([]models.CompaniesTable, error)

	// New method for revenue summary
	GetPaidInvoicesMonthlyRevenue(startDate, endDate *time.Time) ([]struct {
		Month        string
		Year         string
		TotalRevenue float64
	}, error)
}

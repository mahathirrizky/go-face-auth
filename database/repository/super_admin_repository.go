package repository

import (
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type superAdminRepository struct {
	db *gorm.DB
}

func NewSuperAdminRepository(db *gorm.DB) SuperAdminRepository {
	return &superAdminRepository{db: db}
}

// CreateSuperAdmin inserts a new SuperAdmin record into the database.
func (r *superAdminRepository) CreateSuperAdmin(superUser *models.SuperAdminTable) error {
	result := r.db.Create(superUser)
	if result.Error != nil {
		log.Printf("Error creating SuperAdmin: %v", result.Error)
		return result.Error
	}
	log.Printf("SuperAdmin created with ID: %d", superUser.ID)
	return nil
}

// GetSuperAdminByID retrieves a SuperAdmin record by its ID.
func (r *superAdminRepository) GetSuperAdminByID(id int) (*models.SuperAdminTable, error) {
	var superUser models.SuperAdminTable
	result := r.db.First(&superUser, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // SuperAdmin not found
		}
		log.Printf("Error getting SuperAdmin by ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &superUser, nil
}

// GetSuperAdminByEmail retrieves a SuperAdmin record by its email.
func (r *superAdminRepository) GetSuperAdminByEmail(email string) (*models.SuperAdminTable, error) {
	var superUser models.SuperAdminTable
	result := r.db.Where("email = ?", email).First(&superUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // SuperAdmin not found for this email
		}
		log.Printf("Error getting SuperAdmin by email %s: %v", email, result.Error)
		return nil, result.Error
	}
	return &superUser, nil
}

// UpdateSuperAdmin updates an existing SuperAdmin record in the database.
func (r *superAdminRepository) UpdateSuperAdmin(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error) {
	var existingSuperAdmin models.SuperAdminTable
	result := r.db.First(&existingSuperAdmin, id)
	if result.Error != nil {
		log.Printf("Error finding SuperAdmin to update with ID %d: %v", id, result.Error)
		return nil, result.Error
	}

	r.db.Model(&existingSuperAdmin).Updates(superUser)
	return &existingSuperAdmin, nil
}

// DeleteSuperAdmin deletes a SuperAdmin record from the database (soft delete).
func (r *superAdminRepository) DeleteSuperAdmin(id int) error {
	result := r.db.Delete(&models.SuperAdminTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting SuperAdmin with ID %d: %v", id, result.Error)
		return result.Error
	}
	return nil
}

// GetAllSuperAdmins retrieves all SuperAdmin records from the database.
func (r *superAdminRepository) GetAllSuperAdmins() ([]models.SuperAdminTable, error) {
	var superUsers []models.SuperAdminTable
	result := r.db.Find(&superUsers)
	if result.Error != nil {
		log.Printf("Error getting all SuperAdmins: %v", result.Error)
		return nil, result.Error
	}
	return superUsers, nil
}

// GetTotalCompaniesCount returns the total number of companies.
func (r *superAdminRepository) GetTotalCompaniesCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.CompaniesTable{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCompaniesCountBySubscriptionStatus returns the number of companies with a specific subscription status.
func (r *superAdminRepository) GetCompaniesCountBySubscriptionStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.CompaniesTable{}).Where("subscription_status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetExpiredAndTrialExpiredCompaniesCount returns the number of companies with 'expired' or 'expired_trial' status.
func (r *superAdminRepository) GetExpiredAndTrialExpiredCompaniesCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.CompaniesTable{}).Where("subscription_status = ? OR subscription_status = ?", "expired", "expired_trial").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetRecentCompanies returns a limited number of recently created companies.
func (r *superAdminRepository) GetRecentCompanies(limit int) ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	if err := r.db.Order("created_at DESC").Limit(limit).Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// GetAllCompaniesWithPreload returns all companies with their subscription package preloaded.
func (r *superAdminRepository) GetAllCompaniesWithPreload() ([]models.CompaniesTable, error) {
	var companies []models.CompaniesTable
	if err := r.db.Preload("SubscriptionPackage").Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// GetPaidInvoicesMonthlyRevenue returns the monthly revenue from paid invoices within a date range.
func (r *superAdminRepository) GetPaidInvoicesMonthlyRevenue(startDate, endDate *time.Time) ([]struct {
	Month        string
	Year         string
	TotalRevenue float64
}, error) {
	var monthlyRevenue []struct {
		Month        string
		Year         string
		TotalRevenue float64
	}

	query := r.db.Model(&models.InvoiceTable{}).Where("status = ?", "paid")

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	if err := query.Select(
		"DATE_FORMAT(created_at, '%Y-%m') AS month, DATE_FORMAT(created_at, '%Y') AS year, SUM(amount) AS total_revenue").Group("month, year").Order("year DESC, month DESC").Scan(&monthlyRevenue).Error; err != nil {
		return nil, err
	}

	return monthlyRevenue, nil
}

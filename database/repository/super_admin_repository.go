package repository

import (
	"go-face-auth/models"
	"log"

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
package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateSuperAdmin inserts a new SuperAdmin record into the database.
func CreateSuperAdmin(superUser *models.SuperAdminTable) error {
	result := database.DB.Create(superUser)
	if result.Error != nil {
		log.Printf("Error creating SuperAdmin: %v", result.Error)
		return result.Error
	}
	log.Printf("SuperAdmin created with ID: %d", superUser.ID)
	return nil
}

// GetSuperAdminByID retrieves a SuperAdmin record by its ID.
func GetSuperAdminByID(id int) (*models.SuperAdminTable, error) {
	var superUser models.SuperAdminTable
	result := database.DB.First(&superUser, id)
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
func GetSuperAdminByEmail(email string) (*models.SuperAdminTable, error) {
	var superUser models.SuperAdminTable
	result := database.DB.Where("email = ?", email).First(&superUser)
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
func UpdateSuperAdmin(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error) {
	var existingSuperAdmin models.SuperAdminTable
	result := database.DB.First(&existingSuperAdmin, id)
	if result.Error != nil {
		log.Printf("Error finding SuperAdmin to update with ID %d: %v", id, result.Error)
		return nil, result.Error
	}

	database.DB.Model(&existingSuperAdmin).Updates(superUser)
	return &existingSuperAdmin, nil
}

// DeleteSuperAdmin deletes a SuperAdmin record from the database (soft delete).
func DeleteSuperAdmin(id int) error {
	result := database.DB.Delete(&models.SuperAdminTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting SuperAdmin with ID %d: %v", id, result.Error)
		return result.Error
	}
	return nil
}

// GetAllSuperAdmins retrieves all SuperAdmin records from the database.
func GetAllSuperAdmins() ([]models.SuperAdminTable, error) {
	var superUsers []models.SuperAdminTable
	result := database.DB.Find(&superUsers)
	if result.Error != nil {
		log.Printf("Error getting all SuperAdmins: %v", result.Error)
		return nil, result.Error
	}
	return superUsers, nil
}

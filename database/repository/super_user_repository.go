package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateSuperUser inserts a new SuperUser record into the database.
func CreateSuperUser(superUser *models.SuperUserTable) error {
	result := database.DB.Create(superUser)
	if result.Error != nil {
		log.Printf("Error creating SuperUser: %v", result.Error)
		return result.Error
	}
	log.Printf("SuperUser created with ID: %d", superUser.ID)
	return nil
}

// GetSuperUserByID retrieves a SuperUser record by its ID.
func GetSuperUserByID(id int) (*models.SuperUserTable, error) {
	var superUser models.SuperUserTable
	result := database.DB.First(&superUser, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // SuperUser not found
		}
		log.Printf("Error getting SuperUser by ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &superUser, nil
}

// GetSuperUserByEmail retrieves a SuperUser record by its email.
func GetSuperUserByEmail(email string) (*models.SuperUserTable, error) {
	var superUser models.SuperUserTable
	result := database.DB.Where("email = ?", email).First(&superUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // SuperUser not found for this email
		}
		log.Printf("Error getting SuperUser by email %s: %v", email, result.Error)
		return nil, result.Error
	}
	return &superUser, nil
}

// UpdateSuperUser updates an existing SuperUser record in the database.
func UpdateSuperUser(id int, superUser *models.SuperUserTable) (*models.SuperUserTable, error) {
	var existingSuperUser models.SuperUserTable
	result := database.DB.First(&existingSuperUser, id)
	if result.Error != nil {
		log.Printf("Error finding SuperUser to update with ID %d: %v", id, result.Error)
		return nil, result.Error
	}

	database.DB.Model(&existingSuperUser).Updates(superUser)
	return &existingSuperUser, nil
}

// DeleteSuperUser deletes a SuperUser record from the database (soft delete).
func DeleteSuperUser(id int) error {
	result := database.DB.Delete(&models.SuperUserTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting SuperUser with ID %d: %v", id, result.Error)
		return result.Error
	}
	return nil
}

// GetAllSuperUsers retrieves all SuperUser records from the database.
func GetAllSuperUsers() ([]models.SuperUserTable, error) {
	var superUsers []models.SuperUserTable
	result := database.DB.Find(&superUsers)
	if result.Error != nil {
		log.Printf("Error getting all SuperUsers: %v", result.Error)
		return nil, result.Error
	}
	return superUsers, nil
}

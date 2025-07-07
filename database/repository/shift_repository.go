package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
)

// CreateShift inserts a new shift into the database.
func CreateShift(shift *models.ShiftsTable) error {
	result := database.DB.Create(shift)
	if result.Error != nil {
		log.Printf("Error creating shift: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetShiftsByCompanyID retrieves all shifts for a given company ID.
func GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	var shifts []models.ShiftsTable
	result := database.DB.Where("company_id = ?", companyID).Find(&shifts)
	if result.Error != nil {
		log.Printf("Error querying shifts for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return shifts, nil
}

// UpdateShift updates an existing shift record.
func UpdateShift(shift *models.ShiftsTable) error {
	result := database.DB.Save(shift)
	if result.Error != nil {
		log.Printf("Error updating shift: %v", result.Error)
		return result.Error
	}
	return nil
}

// DeleteShift removes a shift from the database by its ID.
func DeleteShift(id int) error {
	result := database.DB.Delete(&models.ShiftsTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting shift with ID %d: %v", id, result.Error)
		return result.Error
	}
	return nil
}

// GetDefaultShiftByCompanyID retrieves the default shift for a given company ID.
func GetDefaultShiftByCompanyID(companyID int) (*models.ShiftsTable, error) {
	var shift models.ShiftsTable
	result := database.DB.Where("company_id = ? AND is_default = ?", companyID, true).First(&shift)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, nil // No default shift found
		}
		log.Printf("Error querying default shift for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return &shift, nil
}

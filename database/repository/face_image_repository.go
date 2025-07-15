package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateFaceImage inserts a new face image into the database.
func CreateFaceImage(faceImage *models.FaceImagesTable) error {
	log.Printf("Attempting to create face image for EmployeeID: %d, ImagePath: %s", faceImage.EmployeeID, faceImage.ImagePath)
	result := database.DB.Create(faceImage)
	if result.Error != nil {
		log.Printf("Error creating face image: %v", result.Error)
		return result.Error
	}
	log.Printf("Face image created with ID: %d, RowsAffected: %d", faceImage.ID, result.RowsAffected)
	return nil
}

// GetFaceImagesByEmployeeID retrieves all face images for a given employee ID.
func GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	var faceImages []models.FaceImagesTable
	log.Printf("Attempting to retrieve face images for EmployeeID: %d", employeeID)
	result := database.DB.Where("employee_id = ?", employeeID).Find(&faceImages)
	if result.Error != nil {
		log.Printf("Error querying face images for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	log.Printf("Found %d face images for EmployeeID: %d", len(faceImages), employeeID)
	return faceImages, nil
}

// GetFaceImageByID retrieves a single face image by its ID.
func GetFaceImageByID(id int) (*models.FaceImagesTable, error) {
	var faceImage models.FaceImagesTable
	result := database.DB.First(&faceImage, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		log.Printf("Error getting face image with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &faceImage, nil
}

// DeleteFaceImage removes a face image record from the database by its ID.
func DeleteFaceImage(id int) error {
	result := database.DB.Delete(&models.FaceImagesTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting face image with ID %d: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("No face image found with ID %d to delete", id)
		return gorm.ErrRecordNotFound // Or a custom error
	}
	log.Printf("Face image with ID %d deleted", id)
	return nil
}

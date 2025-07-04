package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
)

// CreateFaceImage inserts a new face image into the database.
func CreateFaceImage(faceImage *models.FaceImagesTable) error {
	result := database.DB.Create(faceImage)
	if result.Error != nil {
		log.Printf("Error creating face image: %v", result.Error)
		return result.Error
	}
	log.Printf("Face image created with ID: %d", faceImage.ID)
	return nil
}

// GetFaceImagesByEmployeeID retrieves all face images for a given employee ID.
func GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	var faceImages []models.FaceImagesTable
	result := database.DB.Where("employee_id = ?", employeeID).Find(&faceImages)
	if result.Error != nil {
		log.Printf("Error querying face images for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return faceImages, nil
}

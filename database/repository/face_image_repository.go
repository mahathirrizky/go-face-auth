package repository

import (
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

type faceImageRepository struct {
	db *gorm.DB
}

func NewFaceImageRepository(db *gorm.DB) FaceImageRepository {
	return &faceImageRepository{db: db}
}

// CreateFaceImage inserts a new face image into the database.
func (r *faceImageRepository) CreateFaceImage(faceImage *models.FaceImagesTable) error {
	log.Printf("Attempting to create face image for EmployeeID: %d, ImagePath: %s", faceImage.EmployeeID, faceImage.ImagePath)
	result := r.db.Create(faceImage)
	if result.Error != nil {
		log.Printf("Error creating face image: %v", result.Error)
		return result.Error
	}
	log.Printf("Face image created with ID: %d, RowsAffected: %d", faceImage.ID, result.RowsAffected)
	return nil
}

// GetFaceImagesByEmployeeID retrieves all face images for a given employee ID.
func (r *faceImageRepository) GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	var faceImages []models.FaceImagesTable
	log.Printf("Attempting to retrieve face images for EmployeeID: %d", employeeID)
	result := r.db.Where("employee_id = ?", employeeID).Find(&faceImages)
	if result.Error != nil {
		log.Printf("Error querying face images for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	log.Printf("Found %d face images for EmployeeID: %d", len(faceImages), employeeID)
	return faceImages, nil
}

// GetFaceImageByID retrieves a single face image by its ID.
func (r *faceImageRepository) GetFaceImageByID(id int) (*models.FaceImagesTable, error) {
	var faceImage models.FaceImagesTable
	result := r.db.First(&faceImage, id)
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
func (r *faceImageRepository) DeleteFaceImage(id int) error {
	result := r.db.Delete(&models.FaceImagesTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting face image with ID %d: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("No face image found with ID %d to delete or already deleted", id)
		return gorm.ErrRecordNotFound // Or a custom error
	}
	log.Printf("Face image with ID %d deleted", id)
	return nil
}
package repository

import "go-face-auth/models"

// FaceImageRepository defines the contract for face_image-related database operations.
type FaceImageRepository interface {
	CreateFaceImage(faceImage *models.FaceImagesTable) error
	GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error)
	GetFaceImageByID(id int) (*models.FaceImagesTable, error)
	DeleteFaceImage(id int) error
}

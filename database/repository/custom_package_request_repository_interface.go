package repository

import "go-face-auth/models"

// CustomPackageRequestRepository defines the contract for custom_package_request-related database operations.
type CustomPackageRequestRepository interface {
	CreateCustomPackageRequest(req *models.CustomPackageRequest) error
	GetCustomPackageRequestsPaginated(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error)
	GetCustomPackageRequestByID(id uint) (*models.CustomPackageRequest, error)
	UpdateCustomPackageRequest(req *models.CustomPackageRequest) error
}

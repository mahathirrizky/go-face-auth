package repository

import "go-face-auth/models"

// AttendanceLocationRepository defines the contract for attendance_location-related database operations.
type AttendanceLocationRepository interface {
	CreateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error)
	GetAttendanceLocationsByCompanyID(companyID uint) ([]models.AttendanceLocation, error)
	GetAttendanceLocationByID(locationID uint) (*models.AttendanceLocation, error)
	GetLocationsByIDs(ids []uint) ([]models.AttendanceLocation, error)
	UpdateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error)
	DeleteAttendanceLocation(locationID uint) error
	CountAttendanceLocationsByCompanyID(companyID uint) (int64, error)
}

package repository

import (
	"go-face-auth/models"
	"gorm.io/gorm"
)

type attendanceLocationRepository struct {
	db *gorm.DB
}

func NewAttendanceLocationRepository(db *gorm.DB) AttendanceLocationRepository {
	return &attendanceLocationRepository{db: db}
}

func (r *attendanceLocationRepository) CreateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	if err := r.db.Create(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (r *attendanceLocationRepository) GetAttendanceLocationsByCompanyID(companyID uint) ([]models.AttendanceLocation, error) {
	var locations []models.AttendanceLocation
	if err := r.db.Where("company_id = ?", companyID).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *attendanceLocationRepository) GetAttendanceLocationByID(locationID uint) (*models.AttendanceLocation, error) {
	var location models.AttendanceLocation
	if err := r.db.First(&location, locationID).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *attendanceLocationRepository) UpdateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	if err := r.db.Save(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (r *attendanceLocationRepository) DeleteAttendanceLocation(locationID uint) error {
	if err := r.db.Delete(&models.AttendanceLocation{}, locationID).Error; err != nil {
		return err
	}
	return nil
}

func (r *attendanceLocationRepository) CountAttendanceLocationsByCompanyID(companyID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.AttendanceLocation{}).Where("company_id = ?", companyID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
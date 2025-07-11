package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
)

func CreateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	if err := database.DB.Create(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func GetAttendanceLocationsByCompanyID(companyID uint) ([]models.AttendanceLocation, error) {
	var locations []models.AttendanceLocation
	if err := database.DB.Where("company_id = ?", companyID).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func GetAttendanceLocationByID(locationID uint) (*models.AttendanceLocation, error) {
	var location models.AttendanceLocation
	if err := database.DB.First(&location, locationID).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func UpdateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	if err := database.DB.Save(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func DeleteAttendanceLocation(locationID uint) error {
	if err := database.DB.Delete(&models.AttendanceLocation{}, locationID).Error; err != nil {
		return err
	}
	return nil
}

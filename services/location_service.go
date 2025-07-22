package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

func CreateAttendanceLocation(companyID uint, location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").First(&company, companyID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve company information: %w", err)
	}

	var locationCount int64
	if err := database.DB.Model(&models.AttendanceLocation{}).Where("company_id = ?", companyID).Count(&locationCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count existing locations: %w", err)
	}

	if locationCount >= int64(company.SubscriptionPackage.MaxLocations) {
		return nil, fmt.Errorf("location limit reached for your subscription package")
	}

	location.CompanyID = companyID

	return repository.CreateAttendanceLocation(location)
}

func GetAttendanceLocationsByCompanyID(companyID uint) ([]*models.AttendanceLocation, error) {
	locations, err := repository.GetAttendanceLocationsByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	var ptrLocations []*models.AttendanceLocation
	for i := range locations {
		ptrLocations = append(ptrLocations, &locations[i])
	}
	return ptrLocations, nil
}

func UpdateAttendanceLocation(companyID, locationID uint, locationUpdates *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	existingLocation, err := repository.GetAttendanceLocationByID(locationID)
	if err != nil {
		return nil, fmt.Errorf("location not found")
	}
	if existingLocation.CompanyID != companyID {
		return nil, fmt.Errorf("forbidden: you can only update locations for your own company")
	}

	existingLocation.Name = locationUpdates.Name
	existingLocation.Latitude = locationUpdates.Latitude
	existingLocation.Longitude = locationUpdates.Longitude
	existingLocation.Radius = locationUpdates.Radius

	return repository.UpdateAttendanceLocation(existingLocation)
}

func DeleteAttendanceLocation(companyID, locationID uint) error {
	existingLocation, err := repository.GetAttendanceLocationByID(locationID)
	if err != nil {
		return fmt.Errorf("location not found")
	}
	if existingLocation.CompanyID != companyID {
		return fmt.Errorf("forbidden: you can only delete locations for your own company")
	}

	return repository.DeleteAttendanceLocation(locationID)
}

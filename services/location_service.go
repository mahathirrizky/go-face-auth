package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

var ErrLocationLimitReached = fmt.Errorf("location limit reached for your subscription package")

func CreateAttendanceLocation(companyID uint, location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Preload("CustomOffer").First(&company, companyID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve company information: %w", err)
	}

	// Determine the effective MaxLocations limit
	var maxLocationsLimit int
	if company.CustomOfferID != nil && company.CustomOffer != nil {
		maxLocationsLimit = company.CustomOffer.MaxLocations
	} else if company.SubscriptionPackage.ID != 0 {
		maxLocationsLimit = company.SubscriptionPackage.MaxLocations
	} else {
		return nil, fmt.Errorf("company has no active subscription package or custom offer")
	}

	var locationCount int64
	if err := database.DB.Model(&models.AttendanceLocation{}).Where("company_id = ?", companyID).Count(&locationCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count existing locations: %w", err)
	}

	if locationCount >= int64(maxLocationsLimit) {
		return nil, fmt.Errorf("location limit reached for your current plan")
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

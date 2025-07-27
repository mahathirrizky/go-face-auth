package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"

	"gorm.io/gorm"
)

var ErrLocationLimitReached = fmt.Errorf("location limit reached for your subscription package")

// LocationService defines the interface for location related business logic.
type LocationService interface {
	CreateAttendanceLocation(companyID uint, location *models.AttendanceLocation) (*models.AttendanceLocation, error)
	GetAttendanceLocationsByCompanyID(companyID uint) ([]*models.AttendanceLocation, error)
	UpdateAttendanceLocation(companyID, locationID uint, locationUpdates *models.AttendanceLocation) (*models.AttendanceLocation, error)
	DeleteAttendanceLocation(companyID, locationID uint) error
}

// locationService is the concrete implementation of LocationService.
type locationService struct {
	companyRepo          repository.CompanyRepository
	attendanceLocationRepo repository.AttendanceLocationRepository
	db                   *gorm.DB
}

// NewLocationService creates a new instance of LocationService.
func NewLocationService(companyRepo repository.CompanyRepository, attendanceLocationRepo repository.AttendanceLocationRepository, db *gorm.DB) LocationService {
	return &locationService{
		companyRepo:          companyRepo,
		attendanceLocationRepo: attendanceLocationRepo,
		db:                   db,
	}
}

func (s *locationService) CreateAttendanceLocation(companyID uint, location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	company, err := s.companyRepo.GetCompanyWithSubscriptionDetails(int(companyID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve company information: %w", err)
	}

	if company == nil {
		return nil, fmt.Errorf("company with ID %d not found", companyID)
	}

	// Determine the effective MaxLocations limit
	var maxLocationsLimit int
	if company.CustomOfferID != nil && company.CustomOffer != nil {
		maxLocationsLimit = company.CustomOffer.MaxLocations
	} else if company.SubscriptionPackage != nil && company.SubscriptionPackage.ID != 0 {
		maxLocationsLimit = company.SubscriptionPackage.MaxLocations
	} else {
		return nil, fmt.Errorf("company has no active subscription package or custom offer")
	}

	var locationCount int64
	if err := s.db.Model(&models.AttendanceLocation{}).Where("company_id = ?", companyID).Count(&locationCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count existing locations: %w", err)
	}

	if locationCount >= int64(maxLocationsLimit) {
		return nil, ErrLocationLimitReached
	}

	location.CompanyID = companyID

	return s.attendanceLocationRepo.CreateAttendanceLocation(location)
}

func (s *locationService) GetAttendanceLocationsByCompanyID(companyID uint) ([]*models.AttendanceLocation, error) {
	locations, err := s.attendanceLocationRepo.GetAttendanceLocationsByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	var ptrLocations []*models.AttendanceLocation
	for i := range locations {
		ptrLocations = append(ptrLocations, &locations[i])
	}
	return ptrLocations, nil
}

func (s *locationService) UpdateAttendanceLocation(companyID, locationID uint, locationUpdates *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	existingLocation, err := s.attendanceLocationRepo.GetAttendanceLocationByID(locationID)
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

	return s.attendanceLocationRepo.UpdateAttendanceLocation(existingLocation)
}

func (s *locationService) DeleteAttendanceLocation(companyID, locationID uint) error {
	existingLocation, err := s.attendanceLocationRepo.GetAttendanceLocationByID(locationID)
	if err != nil {
		return fmt.Errorf("location not found")
	}
	if existingLocation.CompanyID != companyID {
		return fmt.Errorf("forbidden: you can only delete locations for your own company")
	}

	return s.attendanceLocationRepo.DeleteAttendanceLocation(locationID)
}
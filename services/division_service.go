package services
import (
	"errors"
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"


)

type DivisionService interface {
	CreateDivision(req CreateDivisionRequest) (*models.DivisionTable, error)
	GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error)
	GetDivisionByID(divisionID uint, companyID uint) (*models.DivisionTable, error)
	UpdateDivision(divisionID uint, companyID uint, req UpdateDivisionRequest) (*models.DivisionTable, error)
	DeleteDivision(divisionID uint, companyID uint) error
}

type divisionService struct {
	repo         repository.DivisionRepository
	shiftRepo    repository.ShiftRepository
	locationRepo repository.AttendanceLocationRepository
}

// NewDivisionService creates a new instance of DivisionService.
func NewDivisionService(repo repository.DivisionRepository, shiftRepo repository.ShiftRepository, locationRepo repository.AttendanceLocationRepository) DivisionService {
	return &divisionService{repo: repo, shiftRepo: shiftRepo, locationRepo: locationRepo}
}

// CreateDivisionRequest defines the structure for creating a new division.
type CreateDivisionRequest struct {
	CompanyID   uint   `json:"company_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ShiftIDs    []uint `json:"shift_ids"`
	LocationIDs []uint `json:"location_ids"`
}

func (s *divisionService) CreateDivision(req CreateDivisionRequest) (*models.DivisionTable, error) {
	nameTaken, err := s.repo.IsDivisionNameTaken(req.Name, req.CompanyID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check division name: %w", err)
	}
	if nameTaken {
		return nil, errors.New("division name is already taken")
	}

	var shifts []models.ShiftsTable
	if len(req.ShiftIDs) > 0 {
		shifts, err = s.shiftRepo.GetShiftsByIDs(req.ShiftIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve shifts: %w", err)
		}
		if len(shifts) != len(req.ShiftIDs) {
			return nil, errors.New("one or more selected shifts not found")
		}
	}

	var locations []models.AttendanceLocation
	if len(req.LocationIDs) > 0 {
		locations, err = s.locationRepo.GetLocationsByIDs(req.LocationIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve locations: %w", err)
		}
		if len(locations) != len(req.LocationIDs) {
			return nil, errors.New("one or more selected locations not found")
		}
	}

	division := &models.DivisionTable{
		CompanyID:   req.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		Shifts:      shifts,
		Locations:   locations,
	}
	return s.repo.CreateDivision(division)
}

func (s *divisionService) GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error) {
	return s.repo.GetDivisionsByCompanyID(companyID)
}

func (s *divisionService) GetDivisionByID(divisionID uint, companyID uint) (*models.DivisionTable, error) {
	division, err := s.repo.GetDivisionByID(divisionID)
	if err != nil {
		return nil, err
	}
	if division.CompanyID != companyID {
		return nil, errors.New("division not found in this company")
	}
	return division, nil
}

// UpdateDivisionRequest defines the structure for updating an existing division.
type UpdateDivisionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ShiftIDs    []uint `json:"shift_ids"`
	LocationIDs []uint `json:"location_ids"`
}

func (s *divisionService) UpdateDivision(divisionID uint, companyID uint, req UpdateDivisionRequest) (*models.DivisionTable, error) {
	// Fetch the existing division to ensure it belongs to the company
	existingDivision, err := s.GetDivisionByID(divisionID, companyID)
	if err != nil {
		return nil, err
	}

	// Update fields from the input request only if they are provided
	if req.Name != "" {
		// Check if the new name is already taken by another division in the same company
		nameTaken, err := s.repo.IsDivisionNameTaken(req.Name, companyID, divisionID)
		if err != nil {
			return nil, fmt.Errorf("failed to check division name: %w", err)
		}
		if nameTaken {
			return nil, errors.New("division name is already taken")
		}
		existingDivision.Name = req.Name
	}

	if req.Description != "" {
		existingDivision.Description = req.Description
	}

	// Handle ShiftIDs update
	if req.ShiftIDs != nil {
		var shifts []models.ShiftsTable
		if len(req.ShiftIDs) > 0 {
			shifts, err = s.shiftRepo.GetShiftsByIDs(req.ShiftIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve shifts: %w", err)
			}
			if len(shifts) != len(req.ShiftIDs) {
				return nil, errors.New("one or more selected shifts not found")
			}
		}
		existingDivision.Shifts = shifts
	}

	// Handle LocationIDs update
	if req.LocationIDs != nil {
		var locations []models.AttendanceLocation
		if len(req.LocationIDs) > 0 {
			locations, err = s.locationRepo.GetLocationsByIDs(req.LocationIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve locations: %w", err)
			}
			if len(locations) != len(req.LocationIDs) {
				return nil, errors.New("one or more selected locations not found")
			}
		}
		existingDivision.Locations = locations
	}

	return s.repo.UpdateDivision(existingDivision)
}

func (s *divisionService) DeleteDivision(divisionID uint, companyID uint) error {
	_, err := s.GetDivisionByID(divisionID, companyID)
	if err != nil {
		return err
	}
	return s.repo.DeleteDivision(divisionID)
}

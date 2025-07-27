package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

// ShiftService defines the interface for shift-related business logic.
type ShiftService interface {
	CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error)
	GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error)
	GetShiftByID(id int) (*models.ShiftsTable, error)
	UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error)
	DeleteShift(id int) error
	SetDefaultShift(companyID, shiftID int) error
}

// shiftService is the concrete implementation of ShiftService.
type shiftService struct {
	shiftRepo   repository.ShiftRepository
	companyRepo repository.CompanyRepository
}

// NewShiftService creates a new instance of ShiftService.
func NewShiftService(shiftRepo repository.ShiftRepository, companyRepo repository.CompanyRepository) ShiftService {
	return &shiftService{
		shiftRepo:   shiftRepo,
		companyRepo: companyRepo,
	}
}

func (s *shiftService) CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	company, err := s.companyRepo.GetCompanyWithSubscriptionDetails(shift.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve company information: %w", err)
	}
	if company == nil {
		return nil, fmt.Errorf("company with ID %d not found", shift.CompanyID)
	}

	// Determine the effective MaxShifts limit
	var maxShiftsLimit int
	if company.CustomOfferID != nil && company.CustomOffer != nil {
		maxShiftsLimit = company.CustomOffer.MaxShifts
	} else if company.SubscriptionPackage != nil && company.SubscriptionPackage.ID != 0 {
		maxShiftsLimit = company.SubscriptionPackage.MaxShifts
	} else {
		return nil, fmt.Errorf("company has no active subscription package or custom offer")
	}

	shifts, err := s.shiftRepo.GetShiftsByCompanyID(shift.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing shifts: %w", err)
	}

	if len(shifts) >= maxShiftsLimit {
		return nil, fmt.Errorf("shift limit reached for your current plan")
	}

	return s.shiftRepo.CreateShift(shift)
}

func (s *shiftService) GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	return s.shiftRepo.GetShiftsByCompanyID(companyID)
}

func (s *shiftService) GetShiftByID(id int) (*models.ShiftsTable, error) {
	return s.shiftRepo.GetShiftByID(id)
}

func (s *shiftService) UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	// You might want to add validation here, e.g., check if the shift belongs to the correct company
	return s.shiftRepo.UpdateShift(shift)
}

func (s *shiftService) DeleteShift(id int) error {
	// You might want to add validation here, e.g., check if the shift belongs to the correct company
	return s.shiftRepo.DeleteShift(id)
}

func (s *shiftService) SetDefaultShift(companyID, shiftID int) error {
	// Business logic can be added here, e.g., ensuring the shift exists and belongs to the company
	// For now, we delegate directly to the repository
	return s.shiftRepo.SetDefaultShift(companyID, shiftID)
}
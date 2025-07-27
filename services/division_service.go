package services

import (
	"errors"
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

// DivisionService defines the interface for division-related services.
type DivisionService interface {
	CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error)
	GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error)
	GetDivisionByID(divisionID uint, companyID uint) (*models.DivisionTable, error)
	UpdateDivision(division *models.DivisionTable, companyID uint) (*models.DivisionTable, error)
	DeleteDivision(divisionID uint, companyID uint) error
}

// divisionService is the concrete implementation of DivisionService.
type divisionService struct {
	repo repository.DivisionRepository
}

// NewDivisionService creates a new instance of DivisionService.
func NewDivisionService(repo repository.DivisionRepository) DivisionService {
	return &divisionService{repo: repo}
}

func (s *divisionService) CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	nameTaken, err := s.repo.IsDivisionNameTaken(division.Name, division.CompanyID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check division name: %w", err)
	}
	if nameTaken {
		return nil, errors.New("division name is already taken")
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

func (s *divisionService) UpdateDivision(division *models.DivisionTable, companyID uint) (*models.DivisionTable, error) {
	existingDivision, err := s.GetDivisionByID(division.ID, companyID)
	if err != nil {
		return nil, err
	}

	nameTaken, err := s.repo.IsDivisionNameTaken(division.Name, companyID, division.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check division name: %w", err)
	}
	if nameTaken {
		return nil, errors.New("division name is already taken")
	}

	existingDivision.Name = division.Name
	existingDivision.Description = division.Description

	return s.repo.UpdateDivision(existingDivision)
}

func (s *divisionService) DeleteDivision(divisionID uint, companyID uint) error {
	_, err := s.GetDivisionByID(divisionID, companyID)
	if err != nil {
		return err
	}
	return s.repo.DeleteDivision(divisionID)
}

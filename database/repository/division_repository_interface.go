package repository

import "go-face-auth/models"

// DivisionRepository defines the contract for division-related database operations.
// This allows for mocking in tests and decoupling the service layer from the database implementation.
type DivisionRepository interface {
	CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error)
	GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error)
	GetDivisionByID(divisionID uint) (*models.DivisionTable, error)
	UpdateDivision(division *models.DivisionTable) (*models.DivisionTable, error)
	DeleteDivision(divisionID uint) error
	IsDivisionNameTaken(name string, companyID uint, currentDivisionID uint) (bool, error)
}

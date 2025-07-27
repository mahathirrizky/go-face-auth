package repository

import "go-face-auth/models"

// ShiftRepository defines the contract for shift-related database operations.
type ShiftRepository interface {
	CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error)
	GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error)
	GetShiftByID(id int) (*models.ShiftsTable, error)
	UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error)
	DeleteShift(id int) error
	SetDefaultShift(companyID, shiftID int) error
	GetDefaultShiftByCompanyID(companyID int) (*models.ShiftsTable, error)
}

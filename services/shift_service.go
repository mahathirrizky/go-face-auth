package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
)

func CreateShift(shift *models.ShiftsTable) error {
	shifts, err := repository.GetShiftsByCompanyID(shift.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing shifts: %w", err)
	}

	if len(shifts) >= 4 {
		return fmt.Errorf("maximum of 4 shifts allowed per company")
	}

	if err := repository.CreateShift(shift); err != nil {
		return fmt.Errorf("failed to create shift: %w", err)
	}

	return nil
}

func GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	return repository.GetShiftsByCompanyID(companyID)
}

func UpdateShift(shift *models.ShiftsTable) error {
	return repository.UpdateShift(shift)
}

func DeleteShift(id int) error {
	return repository.DeleteShift(id)
}

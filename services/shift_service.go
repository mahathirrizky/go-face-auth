package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"gorm.io/gorm"
)

var ErrShiftLimitReached = fmt.Errorf("shift limit reached for your subscription package")

func CreateShift(shift *models.ShiftsTable) error {
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Preload("CustomOffer").First(&company, shift.CompanyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("company not found: %w", err)
		}
		return fmt.Errorf("failed to retrieve company information: %w", err)
	}

	// Determine the effective MaxShifts limit
	var maxShiftsLimit int
	if company.CustomOfferID != nil && company.CustomOffer != nil {
		maxShiftsLimit = company.CustomOffer.MaxShifts
	} else if company.SubscriptionPackage.ID != 0 {
		maxShiftsLimit = company.SubscriptionPackage.MaxShifts
	} else {
		return fmt.Errorf("company has no active subscription package or custom offer")
	}

	shifts, err := repository.GetShiftsByCompanyID(shift.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing shifts: %w", err)
	}

	if len(shifts) >= maxShiftsLimit {
		return fmt.Errorf("shift limit reached for your current plan")
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

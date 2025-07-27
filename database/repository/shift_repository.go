package repository

import (
	"go-face-auth/models"

	"gorm.io/gorm"
)

// shiftRepository is the concrete implementation of ShiftRepository
type shiftRepository struct {
	db *gorm.DB
}

// NewShiftRepository creates a new instance of ShiftRepository.
func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	if err := r.db.Create(shift).Error; err != nil {
		return nil, err
	}
	return shift, nil
}

func (r *shiftRepository) GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	var shifts []models.ShiftsTable
	err := r.db.Where("company_id = ?", companyID).Find(&shifts).Error
	return shifts, err
}

func (r *shiftRepository) GetShiftByID(id int) (*models.ShiftsTable, error) {
	var shift models.ShiftsTable
	if err := r.db.First(&shift, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found, without an error
		}
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	if err := r.db.Save(shift).Error; err != nil {
		return nil, err
	}
	return shift, nil
}

func (r *shiftRepository) DeleteShift(id int) error {
	return r.db.Delete(&models.ShiftsTable{}, id).Error
}

// SetDefaultShift handles the transaction to set a new default shift.
func (r *shiftRepository) SetDefaultShift(companyID, shiftID int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Unset the current default shift
	if err := tx.Model(&models.ShiftsTable{}).Where("company_id = ? AND is_default = ?", companyID, true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the new default shift
	if err := tx.Model(&models.ShiftsTable{}).Where("id = ? AND company_id = ?", shiftID, companyID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *shiftRepository) GetDefaultShiftByCompanyID(companyID int) (*models.ShiftsTable, error) {
	var shift models.ShiftsTable
	err := r.db.Where("company_id = ? AND is_default = ?", companyID, true).First(&shift).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No default shift found
		}
		return nil, err
	}
	return &shift, nil
}
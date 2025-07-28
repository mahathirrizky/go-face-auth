package repository

import (

	"go-face-auth/models"
	"gorm.io/gorm"
)

// divisionRepository is the concrete implementation of DivisionRepository
type divisionRepository struct {
	db *gorm.DB
}

// NewDivisionRepository creates a new instance of DivisionRepository.
func NewDivisionRepository(db *gorm.DB) DivisionRepository {
	return &divisionRepository{db: db}
}

func (r *divisionRepository) CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	err := r.db.Create(division).Error
	if err != nil {
		return nil, err
	}
	return division, nil
}

func (r *divisionRepository) GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error) {
	var divisions []models.DivisionTable
	err := r.db.Preload("Shifts").Preload("Locations").Where("company_id = ?", companyID).Find(&divisions).Error
	return divisions, err
}

func (r *divisionRepository) GetDivisionByID(divisionID uint) (*models.DivisionTable, error) {
	var division models.DivisionTable
	err := r.db.Preload("Shifts").Preload("Locations").First(&division, divisionID).Error
	if err != nil {
		return nil, err
	}
	return &division, nil
}

func (r *divisionRepository) UpdateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	// Load the existing division with its associations
	var existingDivision models.DivisionTable
	err := r.db.Preload("Shifts").Preload("Locations").First(&existingDivision, division.ID).Error
	if err != nil {
		return nil, err
	}

	// Update basic fields
	existingDivision.Name = division.Name
	existingDivision.Description = division.Description

	// Replace associations for Shifts
	err = r.db.Model(&existingDivision).Association("Shifts").Replace(division.Shifts)
	if err != nil {
		return nil, err
	}

	// Replace associations for Locations
	err = r.db.Model(&existingDivision).Association("Locations").Replace(division.Locations)
	if err != nil {
		return nil, err
	}

	// Save the updated division (this will also save changes to basic fields)
	err = r.db.Save(&existingDivision).Error
	if err != nil {
		return nil, err
	}

	return &existingDivision, nil
}

func (r *divisionRepository) DeleteDivision(divisionID uint) error {
	return r.db.Delete(&models.DivisionTable{}, divisionID).Error
}

func (r *divisionRepository) IsDivisionNameTaken(name string, companyID uint, currentDivisionID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DivisionTable{}).Where("name = ? AND company_id = ?", name, companyID)
	if currentDivisionID != 0 {
		query = query.Where("id != ?", currentDivisionID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

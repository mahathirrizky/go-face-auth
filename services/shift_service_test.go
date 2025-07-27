package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShift(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	service := services.NewShiftService(mockShiftRepo, mockCompanyRepo)

	shift := &models.ShiftsTable{CompanyID: 1, Name: "Morning", StartTime: "08:00", EndTime: "16:00"}
	company := &models.CompaniesTable{ID: 1, SubscriptionPackage: &models.SubscriptionPackageTable{ID: 1, MaxShifts: 5}}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockShiftRepo.GetShiftsByCompanyIDFunc = func(companyID int) ([]models.ShiftsTable, error) {
			return []models.ShiftsTable{}, nil // No existing shifts
		}
		mockShiftRepo.CreateShiftFunc = func(s *models.ShiftsTable) (*models.ShiftsTable, error) {
			return s, nil
		}

		result, err := service.CreateShift(shift)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shift.Name, result.Name)
	})

	t.Run("Limit Reached", func(t *testing.T) {
		mockCompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockShiftRepo.GetShiftsByCompanyIDFunc = func(companyID int) ([]models.ShiftsTable, error) {
			// Return shifts up to the limit to simulate limit reached
			limitedShifts := make([]models.ShiftsTable, company.SubscriptionPackage.MaxShifts)
			return limitedShifts, nil
		}

		_, err := service.CreateShift(shift)

		assert.Error(t, err)
		assert.Equal(t, "shift limit reached for your current plan", err.Error())
	})
}

func TestGetShiftsByCompanyID(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	service := services.NewShiftService(mockShiftRepo, nil)

	shifts := []models.ShiftsTable{{ID: 1}, {ID: 2}}

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo.GetShiftsByCompanyIDFunc = func(companyID int) ([]models.ShiftsTable, error) {
			return shifts, nil
		}

		result, err := service.GetShiftsByCompanyID(1)

		assert.NoError(t, err)
		assert.Equal(t, shifts, result)
	})
}

func TestGetShiftByID(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	service := services.NewShiftService(mockShiftRepo, nil)

	shift := &models.ShiftsTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo.GetShiftByIDFunc = func(id int) (*models.ShiftsTable, error) {
			return shift, nil
		}

		result, err := service.GetShiftByID(1)

		assert.NoError(t, err)
		assert.Equal(t, shift, result)
	})
}

func TestUpdateShift(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	service := services.NewShiftService(mockShiftRepo, nil)

	shift := &models.ShiftsTable{ID: 1, Name: "Updated Shift"}

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo.UpdateShiftFunc = func(s *models.ShiftsTable) (*models.ShiftsTable, error) {
			return s, nil
		}

		result, err := service.UpdateShift(shift)

		assert.NoError(t, err)
		assert.Equal(t, shift, result)
	})
}

func TestDeleteShift(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	service := services.NewShiftService(mockShiftRepo, nil)

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo.DeleteShiftFunc = func(id int) error {
			return nil
		}

		err := service.DeleteShift(1)

		assert.NoError(t, err)
	})
}

func TestSetDefaultShift(t *testing.T) {
	mockShiftRepo := new(MockShiftRepository)
	service := services.NewShiftService(mockShiftRepo, nil)

	t.Run("Success", func(t *testing.T) {
		mockShiftRepo.SetDefaultShiftFunc = func(companyID, shiftID int) error {
			return nil
		}

		err := service.SetDefaultShift(1, 1)

		assert.NoError(t, err)
	})
}

package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShift(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, mocks.CompanyRepo)

	shift := &models.ShiftsTable{CompanyID: 1, Name: "Morning", StartTime: "08:00", EndTime: "16:00"}
	company := &models.CompaniesTable{ID: 1, SubscriptionPackage: &models.SubscriptionPackageTable{ID: 1, MaxShifts: 5}}

	t.Run("Success", func(t *testing.T) {
		mocks.CompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mocks.ShiftRepo.On("GetShiftsByCompanyID", shift.CompanyID).Return([]models.ShiftsTable{}, nil).Once()
		mocks.ShiftRepo.On("CreateShift", shift).Return(shift, nil).Once()

		result, err := service.CreateShift(shift)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shift.Name, result.Name)
	})

	t.Run("Limit Reached", func(t *testing.T) {
		mocks.CompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		limitedShifts := make([]models.ShiftsTable, company.SubscriptionPackage.MaxShifts)
		mocks.ShiftRepo.On("GetShiftsByCompanyID", shift.CompanyID).Return(limitedShifts, nil).Once()

		_, err := service.CreateShift(shift)

		assert.Error(t, err)
		assert.Equal(t, "shift limit reached for your current plan", err.Error())
	})
}

func TestGetShiftsByCompanyID(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, nil)

	shifts := []models.ShiftsTable{{ID: 1}, {ID: 2}}

	t.Run("Success", func(t *testing.T) {
		mocks.ShiftRepo.On("GetShiftsByCompanyID", 1).Return(shifts, nil).Once()

		result, err := service.GetShiftsByCompanyID(1)

		assert.NoError(t, err)
		assert.Equal(t, shifts, result)
	})
}

func TestGetShiftByID(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, nil)

	shift := &models.ShiftsTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mocks.ShiftRepo.On("GetShiftByID", 1).Return(shift, nil).Once()

		result, err := service.GetShiftByID(1)

		assert.NoError(t, err)
		assert.Equal(t, shift, result)
	})
}

func TestUpdateShift(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, nil)

	shift := &models.ShiftsTable{ID: 1, Name: "Updated Shift"}

	t.Run("Success", func(t *testing.T) {
		mocks.ShiftRepo.On("UpdateShift", shift).Return(shift, nil).Once()

		result, err := service.UpdateShift(shift)

		assert.NoError(t, err)
		assert.Equal(t, shift, result)
	})
}

func TestDeleteShift(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, nil)

	t.Run("Success", func(t *testing.T) {
		mocks.ShiftRepo.On("DeleteShift", 1).Return(nil).Once()

		err := service.DeleteShift(1)

		assert.NoError(t, err)
	})
}

func TestSetDefaultShift(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewShiftService(mocks.ShiftRepo, nil)

	t.Run("Success", func(t *testing.T) {
		mocks.ShiftRepo.On("SetDefaultShift", 1, 1).Return(nil).Once()

		err := service.SetDefaultShift(1, 1)

		assert.NoError(t, err)
	})
}

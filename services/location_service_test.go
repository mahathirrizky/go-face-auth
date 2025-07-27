package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateAttendanceLocation(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockAttendanceLocationRepo := new(MockAttendanceLocationRepository)
	service := services.NewLocationService(mockCompanyRepo, mockAttendanceLocationRepo)

	location := &models.AttendanceLocation{Name: "Office", Latitude: 1.0, Longitude: 1.0, Radius: 100}
	company := &models.CompaniesTable{ID: 1, SubscriptionPackage: &models.SubscriptionPackageTable{ID: 1, MaxLocations: 5}}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockAttendanceLocationRepo.CountAttendanceLocationsByCompanyIDFunc = func(companyID uint) (int64, error) {
			return 0, nil
		}
		mockAttendanceLocationRepo.CreateAttendanceLocationFunc = func(loc *models.AttendanceLocation) (*models.AttendanceLocation, error) {
			return loc, nil
		}

		result, err := service.CreateAttendanceLocation(1, location)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.CompanyID)
	})

	// Note: Testing the actual `db.Model().Where().Count()` chain is complex with simple mocks.
	// This test assumes the count logic is handled correctly or is covered by integration tests.
}

func TestGetAttendanceLocationsByCompanyID(t *testing.T) {
	mockAttendanceLocationRepo := new(MockAttendanceLocationRepository)
	service := services.NewLocationService(nil, mockAttendanceLocationRepo)

	locations := []models.AttendanceLocation{{Name: "Loc1"}, {Name: "Loc2"}}

	t.Run("Success", func(t *testing.T) {
		mockAttendanceLocationRepo.GetAttendanceLocationsByCompanyIDFunc = func(companyID uint) ([]models.AttendanceLocation, error) {
			return locations, nil
		}

		result, err := service.GetAttendanceLocationsByCompanyID(1)

		assert.NoError(t, err)
		assert.Equal(t, len(locations), len(result))
	})
}

func TestUpdateAttendanceLocation(t *testing.T) {
	mockAttendanceLocationRepo := new(MockAttendanceLocationRepository)
	service := services.NewLocationService(nil, mockAttendanceLocationRepo)

	existingLocation := &models.AttendanceLocation{Model: gorm.Model{ID: 1}, CompanyID: 1, Name: "Old Name"}
	updates := &models.AttendanceLocation{Name: "New Name", Latitude: 2.0, Longitude: 2.0, Radius: 200}

	t.Run("Success", func(t *testing.T) {
		mockAttendanceLocationRepo.GetAttendanceLocationByIDFunc = func(id uint) (*models.AttendanceLocation, error) {
			return existingLocation, nil
		}
		mockAttendanceLocationRepo.UpdateAttendanceLocationFunc = func(loc *models.AttendanceLocation) (*models.AttendanceLocation, error) {
			return loc, nil
		}

		result, err := service.UpdateAttendanceLocation(1, 1, updates)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockAttendanceLocationRepo.GetAttendanceLocationByIDFunc = func(id uint) (*models.AttendanceLocation, error) {
			return existingLocation, nil
		}

		_, err := service.UpdateAttendanceLocation(2, 1, updates)

		assert.Error(t, err)
		assert.Equal(t, "forbidden: you can only update locations for your own company", err.Error())
	})
}

func TestDeleteAttendanceLocation(t *testing.T) {
	mockAttendanceLocationRepo := new(MockAttendanceLocationRepository)
	service := services.NewLocationService(nil, mockAttendanceLocationRepo)

	existingLocation := &models.AttendanceLocation{Model: gorm.Model{ID: 1}, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mockAttendanceLocationRepo.GetAttendanceLocationByIDFunc = func(id uint) (*models.AttendanceLocation, error) {
			return existingLocation, nil
		}
		mockAttendanceLocationRepo.DeleteAttendanceLocationFunc = func(id uint) error {
			return nil
		}

		err := service.DeleteAttendanceLocation(1, 1)

		assert.NoError(t, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockAttendanceLocationRepo.GetAttendanceLocationByIDFunc = func(id uint) (*models.AttendanceLocation, error) {
			return existingLocation, nil
		}

		err := service.DeleteAttendanceLocation(2, 1)

		assert.Error(t, err)
		assert.Equal(t, "forbidden: you can only delete locations for your own company", err.Error())
	})
}

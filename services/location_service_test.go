package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateAttendanceLocation(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLocationService(mocks.CompanyRepo, mocks.AttendanceLocationRepo)

	location := &models.AttendanceLocation{Name: "Office", Latitude: 1.0, Longitude: 1.0, Radius: 100}
	company := &models.CompaniesTable{ID: 1, SubscriptionPackage: &models.SubscriptionPackageTable{ID: 1, MaxLocations: 5}}

	t.Run("Success", func(t *testing.T) {
		mocks.CompanyRepo.GetCompanyWithSubscriptionDetailsFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mocks.AttendanceLocationRepo.On("CountAttendanceLocationsByCompanyID", uint(1)).Return(int64(0), nil).Once()
		mocks.AttendanceLocationRepo.On("CreateAttendanceLocation", location).Return(location, nil).Once()

		result, err := service.CreateAttendanceLocation(1, location)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.CompanyID)
	})

	// Note: Testing the actual `db.Model().Where().Count()` chain is complex with simple mocks.
	// This test assumes the count logic is handled correctly or is covered by integration tests.
}

func TestGetAttendanceLocationsByCompanyID(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLocationService(nil, mocks.AttendanceLocationRepo)

	locations := []models.AttendanceLocation{{Name: "Loc1"}, {Name: "Loc2"}}

	t.Run("Success", func(t *testing.T) {
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationsByCompanyID", uint(1)).Return(locations, nil).Once()

		result, err := service.GetAttendanceLocationsByCompanyID(1)

		assert.NoError(t, err)
		assert.Equal(t, len(locations), len(result))
	})
}

func TestUpdateAttendanceLocation(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLocationService(nil, mocks.AttendanceLocationRepo)

	existingLocation := &models.AttendanceLocation{Model: gorm.Model{ID: 1}, CompanyID: 1, Name: "Old Name"}
	updates := &models.AttendanceLocation{Name: "New Name", Latitude: 2.0, Longitude: 2.0, Radius: 200}

	t.Run("Success", func(t *testing.T) {
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationByID", uint(1)).Return(existingLocation, nil).Once()
		updates.ID = existingLocation.ID
		updates.CompanyID = existingLocation.CompanyID
		mocks.AttendanceLocationRepo.On("UpdateAttendanceLocation", updates).Return(updates, nil).Once()

		result, err := service.UpdateAttendanceLocation(1, 1, updates)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationByID", uint(1)).Return(existingLocation, nil).Once()



				_, err := service.UpdateAttendanceLocation(2, 1, updates)

		assert.Error(t, err)
		assert.Equal(t, "forbidden: you can only update locations for your own company", err.Error())
	})
}

func TestDeleteAttendanceLocation(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLocationService(nil, mocks.AttendanceLocationRepo)

	existingLocation := &models.AttendanceLocation{Model: gorm.Model{ID: 1}, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationByID", uint(1)).Return(existingLocation, nil).Once()
		mocks.AttendanceLocationRepo.On("DeleteAttendanceLocation", uint(1)).Return(nil).Once()

		err := service.DeleteAttendanceLocation(1, 1)

		assert.NoError(t, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationByID", uint(1)).Return(existingLocation, nil).Once()

		err := service.DeleteAttendanceLocation(2, 1)

		assert.Error(t, err)
		assert.Equal(t, "forbidden: you can only delete locations for your own company", err.Error())
	})
}

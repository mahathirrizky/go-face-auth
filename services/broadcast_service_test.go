package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestCreateBroadcastMessage(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewBroadcastService(mocks.BroadcastRepo)

	companyID := uint(1)
	message := "test message"
	expireDate := "2025-12-31"

	t.Run("Success", func(t *testing.T) {
		mocks.BroadcastRepo.CreateBroadcastFunc = func(msg *models.BroadcastMessage) error {
			return nil
		}

		result, err := service.CreateBroadcastMessage(companyID, message, expireDate)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, companyID, result.CompanyID)
		assert.Equal(t, message, result.Message)
	})

	t.Run("Invalid Date", func(t *testing.T) {
		_, err := service.CreateBroadcastMessage(companyID, message, "invalid-date")

		assert.Error(t, err)
	})

	t.Run("Error Creating", func(t *testing.T) {
		mocks.BroadcastRepo.CreateBroadcastFunc = func(msg *models.BroadcastMessage) error {
			return errors.New("db error")
		}

		_, err := service.CreateBroadcastMessage(companyID, message, expireDate)

		assert.Error(t, err)
	})
}

func TestGetBroadcastsForEmployee(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewBroadcastService(mocks.BroadcastRepo)

	companyID := uint(1)
	employeeID := uint(1)
	broadcasts := []models.BroadcastMessage{{Message: "msg1"}, {Message: "msg2"}}

	t.Run("Success", func(t *testing.T) {
		mocks.BroadcastRepo.GetBroadcastsForEmployeeFunc = func(cID, eID uint) ([]models.BroadcastMessage, error) {
			return broadcasts, nil
		}

		result, err := service.GetBroadcastsForEmployee(companyID, employeeID)

		assert.NoError(t, err)
		assert.Equal(t, broadcasts, result)
	})

	t.Run("Error", func(t *testing.T) {
		mocks.BroadcastRepo.GetBroadcastsForEmployeeFunc = func(cID, eID uint) ([]models.BroadcastMessage, error) {
			return nil, errors.New("db error")
		}

		_, err := service.GetBroadcastsForEmployee(companyID, employeeID)

		assert.Error(t, err)
	})
}

func TestMarkBroadcastAsRead(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewBroadcastService(mocks.BroadcastRepo)

	employeeID := uint(1)
	messageID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mocks.BroadcastRepo.MarkBroadcastAsReadFunc = func(eID, mID uint) error {
			return nil
		}

		err := service.MarkBroadcastAsRead(employeeID, messageID)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mocks.BroadcastRepo.MarkBroadcastAsReadFunc = func(eID, mID uint) error {
			return errors.New("db error")
		}

		err := service.MarkBroadcastAsRead(employeeID, messageID)

		assert.Error(t, err)
	})
}

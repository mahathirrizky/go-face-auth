package services_test

import (

	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubscriptionPackages(t *testing.T) {
	mockSubscriptionPackageRepo := new(MockSubscriptionPackageRepository)
	service := services.NewSubscriptionPackageService(mockSubscriptionPackageRepo)

	packages := []models.SubscriptionPackageTable{{ID: 1}, {ID: 2}}

	t.Run("Success", func(t *testing.T) {
		mockSubscriptionPackageRepo.GetSubscriptionPackagesFunc = func() ([]models.SubscriptionPackageTable, error) {
			return packages, nil
		}

		result, err := service.GetSubscriptionPackages()

		assert.NoError(t, err)
		assert.Equal(t, packages, result)
	})
}

func TestCreateSubscriptionPackage(t *testing.T) {
	mockSubscriptionPackageRepo := new(MockSubscriptionPackageRepository)
	service := services.NewSubscriptionPackageService(mockSubscriptionPackageRepo)

	pkg := &models.SubscriptionPackageTable{PackageName: "Basic"}

	t.Run("Success", func(t *testing.T) {
		mockSubscriptionPackageRepo.CreateSubscriptionPackageFunc = func(p *models.SubscriptionPackageTable) error {
			return nil
		}

		err := service.CreateSubscriptionPackage(pkg)

		assert.NoError(t, err)
	})
}

func TestUpdateSubscriptionPackage(t *testing.T) {
	mockSubscriptionPackageRepo := new(MockSubscriptionPackageRepository)
	service := services.NewSubscriptionPackageService(mockSubscriptionPackageRepo)

	updates := map[string]interface{}{"package_name": "Updated Basic"}

	t.Run("Success", func(t *testing.T) {
		mockSubscriptionPackageRepo.UpdateSubscriptionPackageFieldsFunc = func(id string, u map[string]interface{}) error {
			return nil
		}

		err := service.UpdateSubscriptionPackage("1", updates)

		assert.NoError(t, err)
	})
}

func TestDeleteSubscriptionPackage(t *testing.T) {
	mockSubscriptionPackageRepo := new(MockSubscriptionPackageRepository)
	service := services.NewSubscriptionPackageService(mockSubscriptionPackageRepo)

	t.Run("Success", func(t *testing.T) {
		mockSubscriptionPackageRepo.DeleteSubscriptionPackageFunc = func(id string) error {
			return nil
		}

		err := service.DeleteSubscriptionPackage("1")

		assert.NoError(t, err)
	})
}

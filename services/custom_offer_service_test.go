package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomOffer(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewCustomOfferService(mocks.CustomOfferRepo)

	t.Run("Success", func(t *testing.T) {
		mocks.CustomOfferRepo.CreateCustomOfferFunc = func(offer *models.CustomOffer) error {
			return nil
		}

		offer, err := service.CreateCustomOffer(1, "Test Co", "Premium", 100, 10, 5, "All", 1000, "monthly")

		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, "pending", offer.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mocks.CustomOfferRepo.CreateCustomOfferFunc = func(offer *models.CustomOffer) error {
			return errors.New("db error")
		}

		_, err := service.CreateCustomOffer(1, "Test Co", "Premium", 100, 10, 5, "All", 1000, "monthly")

		assert.Error(t, err)
	})
}

func TestGetCustomOfferByToken(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewCustomOfferService(mocks.CustomOfferRepo)

	offer := &models.CustomOffer{Token: "test-token", CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mocks.CustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}

		result, err := service.GetCustomOfferByToken("test-token", 1)

		assert.NoError(t, err)
		assert.Equal(t, offer, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mocks.CustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}

		_, err := service.GetCustomOfferByToken("test-token", 2)

		assert.Error(t, err)
		assert.Equal(t, "unauthorized access to custom offer", err.Error())
	})
}

func TestMarkCustomOfferAsUsed(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewCustomOfferService(mocks.CustomOfferRepo)

	offer := &models.CustomOffer{Token: "test-token", Status: "pending"}

	t.Run("Success", func(t *testing.T) {
		mocks.CustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}
		mocks.CustomOfferRepo.UpdateCustomOfferFunc = func(o *models.CustomOffer) error {
			return nil
		}

		err := service.MarkCustomOfferAsUsed("test-token")

		assert.NoError(t, err)
		assert.Equal(t, "used", offer.Status)
	})
}

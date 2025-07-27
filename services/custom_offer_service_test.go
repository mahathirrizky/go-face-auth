package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomOffer(t *testing.T) {
	mockCustomOfferRepo := new(MockCustomOfferRepository)
	service := services.NewCustomOfferService(mockCustomOfferRepo)

	t.Run("Success", func(t *testing.T) {
		mockCustomOfferRepo.CreateCustomOfferFunc = func(offer *models.CustomOffer) error {
			return nil
		}

		offer, err := service.CreateCustomOffer(1, "Test Co", "Premium", 100, 10, 5, "All", 1000, "monthly")

		assert.NoError(t, err)
		assert.NotNil(t, offer)
		assert.Equal(t, "pending", offer.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockCustomOfferRepo.CreateCustomOfferFunc = func(offer *models.CustomOffer) error {
			return errors.New("db error")
		}

		_, err := service.CreateCustomOffer(1, "Test Co", "Premium", 100, 10, 5, "All", 1000, "monthly")

		assert.Error(t, err)
	})
}

func TestGetCustomOfferByToken(t *testing.T) {
	mockCustomOfferRepo := new(MockCustomOfferRepository)
	service := services.NewCustomOfferService(mockCustomOfferRepo)

	offer := &models.CustomOffer{Token: "test-token", CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mockCustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}

		result, err := service.GetCustomOfferByToken("test-token", 1)

		assert.NoError(t, err)
		assert.Equal(t, offer, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockCustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}

		_, err := service.GetCustomOfferByToken("test-token", 2)

		assert.Error(t, err)
		assert.Equal(t, "unauthorized access to custom offer", err.Error())
	})
}

func TestMarkCustomOfferAsUsed(t *testing.T) {
	mockCustomOfferRepo := new(MockCustomOfferRepository)
	service := services.NewCustomOfferService(mockCustomOfferRepo)

	offer := &models.CustomOffer{Token: "test-token", Status: "pending"}

	t.Run("Success", func(t *testing.T) {
		mockCustomOfferRepo.GetCustomOfferByTokenFunc = func(token string) (*models.CustomOffer, error) {
			return offer, nil
		}
		mockCustomOfferRepo.UpdateCustomOfferFunc = func(o *models.CustomOffer) error {
			return nil
		}

		err := service.MarkCustomOfferAsUsed("test-token")

		assert.NoError(t, err)
		assert.Equal(t, "used", offer.Status)
	})
}

package repository

import "go-face-auth/models"

// CustomOfferRepository defines the contract for custom_offer-related database operations.
type CustomOfferRepository interface {
	CreateCustomOffer(offer *models.CustomOffer) error
	GetCustomOfferByToken(token string) (*models.CustomOffer, error)
	GetCustomOfferByID(id uint) (*models.CustomOffer, error)
	UpdateCustomOffer(offer *models.CustomOffer) error
	MarkCustomOfferAsUsed(token string) error
}

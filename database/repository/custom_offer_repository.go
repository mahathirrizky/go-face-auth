package repository

import (
	"fmt"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

type customOfferRepository struct {
	db *gorm.DB
}

func NewCustomOfferRepository(db *gorm.DB) CustomOfferRepository {
	return &customOfferRepository{db: db}
}

// CreateCustomOffer inserts a new custom offer into the database.
func (r *customOfferRepository) CreateCustomOffer(offer *models.CustomOffer) error {
	result := r.db.Create(offer)
	if result.Error != nil {
		log.Printf("Error creating custom offer: %v", result.Error)
		return result.Error
	}
	log.Printf("Custom offer created with ID: %d, Token: %s", offer.ID, offer.Token)
	return nil
}

// GetCustomOfferByToken retrieves a custom offer by its unique token.
func (r *customOfferRepository) GetCustomOfferByToken(token string) (*models.CustomOffer, error) {
	var offer models.CustomOffer
	result := r.db.Where("token = ?", token).First(&offer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		log.Printf("Error getting custom offer by token %s: %v", token, result.Error)
		return nil, result.Error
	}
	return &offer, nil
}

// GetCustomOfferByID retrieves a custom offer by its ID.
func (r *customOfferRepository) GetCustomOfferByID(id uint) (*models.CustomOffer, error) {
	var offer models.CustomOffer
	result := r.db.First(&offer, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		log.Printf("Error getting custom offer by ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &offer, nil
}

// UpdateCustomOffer updates an existing custom offer.
func (r *customOfferRepository) UpdateCustomOffer(offer *models.CustomOffer) error {
	result := r.db.Save(offer)
	if result.Error != nil {
		log.Printf("Error updating custom offer %d: %v", offer.ID, result.Error)
		return result.Error
	}
	log.Printf("Custom offer %d updated successfully.", offer.ID)
	return nil
}

// MarkCustomOfferAsUsed updates the status of a custom offer to 'used'.
func (r *customOfferRepository) MarkCustomOfferAsUsed(token string) error {
	offer, err := r.GetCustomOfferByToken(token)
	if err != nil {
		return fmt.Errorf("failed to find custom offer: %w", err)
	}
	if offer == nil {
		return fmt.Errorf("custom offer not found")
	}

	offer.Status = "used"
	if err := r.UpdateCustomOffer(offer); err != nil {
		return fmt.Errorf("failed to mark custom offer as used: %w", err)
	}
	return nil
}
package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateCustomOffer inserts a new custom offer into the database.
func CreateCustomOffer(offer *models.CustomOffer) error {
	result := database.DB.Create(offer)
	if result.Error != nil {
		log.Printf("Error creating custom offer: %v", result.Error)
		return result.Error
	}
	log.Printf("Custom offer created with ID: %d, Token: %s", offer.ID, offer.Token)
	return nil
}

// GetCustomOfferByToken retrieves a custom offer by its unique token.
func GetCustomOfferByToken(token string) (*models.CustomOffer, error) {
	var offer models.CustomOffer
	result := database.DB.Where("token = ?", token).First(&offer)
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
func GetCustomOfferByID(id uint) (*models.CustomOffer, error) {
	var offer models.CustomOffer
	result := database.DB.First(&offer, id)
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
func UpdateCustomOffer(offer *models.CustomOffer) error {
	result := database.DB.Save(offer)
	if result.Error != nil {
		log.Printf("Error updating custom offer %d: %v", offer.ID, result.Error)
		return result.Error
	}
	log.Printf("Custom offer %d updated successfully.", offer.ID)
	return nil
}

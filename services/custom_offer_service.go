package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"go-face-auth/helper"
)

// CustomOfferService defines the interface for custom offer related business logic.
type CustomOfferService interface {
	CreateCustomOffer(companyID uint, companyName, packageName string, maxEmployees int, maxLocations int, maxShifts int, features string, finalPrice float64, billingCycle string) (*models.CustomOffer, error)
	GetCustomOfferByToken(token string, companyID uint) (*models.CustomOffer, error)
	MarkCustomOfferAsUsed(token string) error
}

// customOfferService is the concrete implementation of CustomOfferService.
type customOfferService struct {
	customOfferRepo repository.CustomOfferRepository
}

// NewCustomOfferService creates a new instance of CustomOfferService.
func NewCustomOfferService(customOfferRepo repository.CustomOfferRepository) CustomOfferService {
	return &customOfferService{
		customOfferRepo: customOfferRepo,
	}
}

// CreateCustomOffer generates a new custom offer and saves it to the database.
func (s *customOfferService) CreateCustomOffer(companyID uint, companyName, packageName string, maxEmployees int, maxLocations int, maxShifts int, features string, finalPrice float64, billingCycle string) (*models.CustomOffer, error) {
	// Generate a unique token for the offer
	token, err := helper.GenerateRandomString(32) // Using a helper for random string generation
	if err != nil {
		return nil, fmt.Errorf("failed to generate offer token: %w", err)
	}

	offer := &models.CustomOffer{
		Token:        token,
		CompanyID:    companyID,
		CompanyName:  companyName,
		PackageName:  packageName,
		MaxEmployees: maxEmployees,
		MaxLocations: maxLocations,
		MaxShifts:    maxShifts,
		Features:     features,
		FinalPrice:   finalPrice,
		BillingCycle: billingCycle,
		Status:       "pending", // Initial status
	}

	if err := s.customOfferRepo.CreateCustomOffer(offer); err != nil {
		return nil, fmt.Errorf("failed to create custom offer: %w", err)
	}

	return offer, nil
}

// GetCustomOfferByToken retrieves a custom offer by its token.
func (s *customOfferService) GetCustomOfferByToken(token string, companyID uint) (*models.CustomOffer, error) {
	offer, err := s.customOfferRepo.GetCustomOfferByToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve custom offer: %w", err)
	}
	if offer == nil {
		return nil, fmt.Errorf("custom offer not found")
	}

	// Verify that the offer belongs to the requesting company
	if offer.CompanyID != companyID {
		return nil, fmt.Errorf("unauthorized access to custom offer")
	}

	// Optional: Add logic here to check for offer expiration if needed

	return offer, nil
}

// MarkCustomOfferAsUsed updates the status of a custom offer to 'used'.
func (s *customOfferService) MarkCustomOfferAsUsed(token string) error {
	offer, err := s.customOfferRepo.GetCustomOfferByToken(token)
	if err != nil {
		return fmt.Errorf("failed to find custom offer: %w", err)
	}
	if offer == nil {
		return fmt.Errorf("custom offer not found")
	}

	offer.Status = "used"
	if err := s.customOfferRepo.UpdateCustomOffer(offer); err != nil {
		return fmt.Errorf("failed to mark custom offer as used: %w", err)
	}
	return nil
}
package handlers

import (
	"fmt"
	"go-face-auth/helper"
	"go-face-auth/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomOfferHandler defines the interface for custom offer related handlers.
type CustomOfferHandler interface {
	HandleCreateCustomOffer(c *gin.Context)
	HandleGetCustomOfferByToken(c *gin.Context)
}

// customOfferHandler is the concrete implementation of CustomOfferHandler.
type customOfferHandler struct {
	customOfferService services.CustomOfferService
}

// NewCustomOfferHandler creates a new instance of CustomOfferHandler.
func NewCustomOfferHandler(customOfferService services.CustomOfferService) CustomOfferHandler {
	return &customOfferHandler{
		customOfferService: customOfferService,
	}
}

// CreateCustomOfferRequest represents the request body for creating a custom offer.
type CreateCustomOfferRequest struct {
	CompanyID    uint    `json:"company_id" binding:"required"`
	CompanyName  string  `json:"company_name" binding:"required"`
	PackageName  string  `json:"package_name" binding:"required"`
	MaxEmployees int     `json:"max_employees" binding:"required,min=1"`
	MaxLocations int     `json:"max_locations" binding:"required,min=0"`
	MaxShifts    int     `json:"max_shifts" binding:"required,min=0"`
	Features     string  `json:"features" binding:"required"`
	FinalPrice   float64 `json:"final_price" binding:"required,min=0"`
	BillingCycle string  `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
}

// HandleCreateCustomOffer handles the creation of a new custom offer.
func (h *customOfferHandler) HandleCreateCustomOffer(c *gin.Context) {
	var req CreateCustomOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
		return
	}

	offer, err := h.customOfferService.CreateCustomOffer(
		req.CompanyID,
		req.CompanyName,
		req.PackageName,
		req.MaxEmployees,
		req.MaxLocations,
		req.MaxShifts,
		req.Features,
		req.FinalPrice,
		req.BillingCycle,
	)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Construct the offer link (assuming frontend will handle the /offer/:token route)
	offerLink := fmt.Sprintf("%s/offer/%s", helper.GetFrontendAdminURL(), offer.Token)

	helper.SendSuccess(c, http.StatusCreated, "Custom offer created successfully.", gin.H{
		"offer_id": offer.ID,
		"token":    offer.Token,
		"link":     offerLink,
	})
}

// HandleGetCustomOfferByToken retrieves a custom offer by its token.
func (h *customOfferHandler) HandleGetCustomOfferByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		helper.SendError(c, http.StatusBadRequest, "Offer token is required.")
		return
	}

	companyIDClaim, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}
	compID := uint(companyIDClaim.(float64))

	offer, err := h.customOfferService.GetCustomOfferByToken(token, compID)
	if err != nil {
		// Differentiate between not found and unauthorized
		if err.Error() == "custom offer not found" {
			helper.SendError(c, http.StatusNotFound, err.Error())
		} else if err.Error() == "unauthorized access to custom offer" {
			helper.SendError(c, http.StatusForbidden, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Custom offer retrieved successfully.", offer)
}
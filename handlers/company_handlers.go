package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

// UpdateCompanyRequest represents the request body for updating company details.
type UpdateCompanyRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Timezone string `json:"timezone"`
}

// GetCompanyDetails handles fetching company details for the authenticated admin.
func GetCompanyDetails(c *gin.Context) {
	// Get companyID from JWT claims set by AuthMiddleware
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}

	// Ensure companyID is of the correct type (int)
	id, ok := companyID.(float64) // JWT claims are typically float64
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	company, err := repository.GetCompanyByID(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company details.")
		return
	}

	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	// Fetch admin company details using the company ID
	adminCompany, err := repository.GetAdminCompanyByCompanyID(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve admin details for company.")
		return
	}

	// Prepare response data including admin email and timezone
	responseData := gin.H{
		"id":                    company.ID,
		"name":                  company.Name,
		"address":               company.Address,
		"admin_email":           adminCompany.Email, // Include admin email
		"subscription_status":   company.SubscriptionStatus,
		"trial_start_date":      company.TrialStartDate,
		"trial_end_date":        company.TrialEndDate,
		"timezone":              company.Timezone, // Include timezone
	}

	helper.SendSuccess(c, http.StatusOK, "Company details fetched successfully.", responseData)
}

// UpdateCompanyDetails handles updating company details for the authenticated admin.
func UpdateCompanyDetails(c *gin.Context) {
	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	id, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	company, err := repository.GetCompanyByID(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company.")
		return
	}
	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	// Update fields if provided
	if req.Name != "" {
		company.Name = req.Name
	}
	if req.Address != "" {
		company.Address = req.Address
	}
	if req.Timezone != "" {
		// Validate timezone string (optional but recommended)
		_, err := time.LoadLocation(req.Timezone)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid timezone string.")
			return
		}
		company.Timezone = req.Timezone
	}

	if err := repository.UpdateCompany(company); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update company details.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company details updated successfully.", company)
}

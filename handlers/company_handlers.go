package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
)

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

	// Prepare response data including admin email
	responseData := gin.H{
		"id":                    company.ID,
		"name":                  company.Name,
		"address":               company.Address,
		"admin_email":           adminCompany.Email, // Include admin email
		"subscription_status":   company.SubscriptionStatus,
		"trial_start_date":      company.TrialStartDate,
		"trial_end_date":        company.TrialEndDate,
	}

	helper.SendSuccess(c, http.StatusOK, "Company details fetched successfully.", responseData)
}

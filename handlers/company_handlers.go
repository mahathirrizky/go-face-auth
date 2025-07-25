package handlers

import (
	"go-face-auth/services"
	"log"
	"net/http"

	"go-face-auth/websocket"

	"go-face-auth/helper"

	"github.com/gin-gonic/gin"
)

// UpdateCompanyRequest represents the request body for updating company details.
type UpdateCompanyRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Timezone string `json:"timezone"`
}

// RegisterCompanyRequest defines the structure for the company registration request body.
type RegisterCompanyRequest struct {
	CompanyName         string `json:"company_name" binding:"required"`
	CompanyAddress      string `json:"company_address"`
	AdminEmail          string `json:"admin_email" binding:"required,email"`
	AdminPassword       string `json:"admin_password" binding:"required,min=6"`
	SubscriptionPackageID int   `json:"subscription_package_id" binding:"required"`
	BillingCycle        string `json:"billing_cycle" binding:"required"`
}

// RegisterCompany handles the registration of a new company and its admin user.
func RegisterCompany(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterCompanyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err) // Backend Log
			helper.SendError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		log.Printf("Received registration request: %+v", req) // Backend Log

		company, adminCompany, err := services.RegisterCompany(services.RegisterCompanyRequest(req))
		if err != nil {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Trigger a dashboard update
		go hub.BroadcastSuperAdminDashboardUpdate()

		helper.SendSuccess(c, http.StatusCreated, "Company registered successfully. Please check your email for a confirmation link.", gin.H{
			"company_id": company.ID,
			"admin_email": adminCompany.Email,
		})
	}
}


func ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.SendError(c, http.StatusBadRequest, "Confirmation token is missing.")
		return
	}

	if err := services.ConfirmEmail(token); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Email confirmed successfully. You can now log in.", nil)
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

	responseData, err := services.GetCompanyDetails(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company details.")
		return
	}
	if responseData == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company details fetched successfully.", responseData)
}



// GetCompanySubscriptionStatus handles fetching the subscription status for the authenticated company.
func GetCompanySubscriptionStatus(c *gin.Context) {
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

	subscriptionStatus, err := services.GetCompanySubscriptionStatus(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company subscription status fetched successfully.", subscriptionStatus)
}

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

	company, err := services.UpdateCompanyDetails(int(id), req.Name, req.Address, req.Timezone)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update company details.")
		return
	}
	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company details updated successfully.", company)
}

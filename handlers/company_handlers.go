package handlers

import (
	"go-face-auth/services"
	"log"
	"net/http"
	"strconv"

	"go-face-auth/websocket"

	"go-face-auth/helper"

	"github.com/gin-gonic/gin"
)

// CompanyHandler defines the interface for company related handlers.
type CompanyHandler interface {
	RegisterCompany(hub *websocket.Hub) gin.HandlerFunc
	ConfirmEmail(c *gin.Context)
	GetCompanyDetails(c *gin.Context)
	GetCompanySubscriptionStatus(c *gin.Context)
	UpdateCompanyDetails(c *gin.Context)
	CreateCompany(c *gin.Context)
	GetCompanyByID(c *gin.Context)
}

// companyHandler is the concrete implementation of CompanyHandler.
type companyHandler struct {
	companyService services.CompanyService
}

// NewCompanyHandler creates a new instance of CompanyHandler.
func NewCompanyHandler(companyService services.CompanyService) CompanyHandler {
	return &companyHandler{
		companyService: companyService,
	}
}

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
func (h *companyHandler) RegisterCompany(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterCompanyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err) // Backend Log
			helper.SendError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		log.Printf("Received registration request: %+v", req) // Backend Log

		company, adminCompany, err := h.companyService.RegisterCompany(services.RegisterCompanyRequest(req))
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


func (h *companyHandler) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.SendError(c, http.StatusBadRequest, "Confirmation token is missing.")
		return
	}

	if err := h.companyService.ConfirmEmail(token); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Email confirmed successfully. You can now log in.", nil)
}


// GetCompanyDetails handles fetching company details for the authenticated admin.
func (h *companyHandler) GetCompanyDetails(c *gin.Context) {
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

	responseData, err := h.companyService.GetCompanyDetails(int(id))
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
func (h *companyHandler) GetCompanySubscriptionStatus(c *gin.Context) {
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

	subscriptionStatus, err := h.companyService.GetCompanySubscriptionStatus(int(id))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company subscription status fetched successfully.", subscriptionStatus)
}

func (h *companyHandler) UpdateCompanyDetails(c *gin.Context) {
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

	company, err := h.companyService.UpdateCompanyDetails(int(id), req.Name, req.Address, req.Timezone)
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

// CreateCompanyRequest defines the structure for the company creation request body.
type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

// CreateCompany handles the creation of a new company.
func (h *companyHandler) CreateCompany(c *gin.Context) {
	var req services.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	company, err := h.companyService.CreateCompany(req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Company created successfully.", company)
}

// GetCompanyByID handles fetching a company by its ID.
func (h *companyHandler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company.")
		return
	}

	if company == nil {
		helper.SendError(c, http.StatusNotFound, "Company not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Company retrieved successfully.", company)
}

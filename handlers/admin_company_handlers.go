package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/services"
	"go-face-auth/websocket"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminCompanyHandler defines the interface for admin company related handlers.
type AdminCompanyHandler interface {
	CreateAdminCompany(c *gin.Context)
	GetAdminCompanyByCompanyID(c *gin.Context)
	GetAdminCompanyByEmployeeID(c *gin.Context)
	ChangeAdminPassword(c *gin.Context)
	CheckAndNotifySubscriptions(c *gin.Context)
	GetDashboardSummary(hub *websocket.Hub, c *gin.Context)
}

// adminCompanyHandler is the concrete implementation of AdminCompanyHandler.
type adminCompanyHandler struct {
	adminCompanyService services.AdminCompanyService
}

// NewAdminCompanyHandler creates a new instance of AdminCompanyHandler.
func NewAdminCompanyHandler(adminCompanyService services.AdminCompanyService) AdminCompanyHandler {
	return &adminCompanyHandler{
		adminCompanyService: adminCompanyService,
	}
}

// ChangePasswordRequest defines the structure for the change password request body.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// CreateAdminCompany handles the creation of a new admin company.
func (h *adminCompanyHandler) CreateAdminCompany(c *gin.Context) {
	var adminCompany models.AdminCompaniesTable
	if err := c.BindJSON(&adminCompany); err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.adminCompanyService.CreateAdminCompany(&adminCompany); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Admin company created successfully", adminCompany)
}

// GetAdminCompanyByCompanyID handles fetching an admin company by its CompanyID.
func (h *adminCompanyHandler) GetAdminCompanyByCompanyID(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID")
		return
	}

	adminCompany, err := h.adminCompanyService.GetAdminCompanyByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this company ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}

// GetAdminCompanyByEmployeeID handles fetching an admin company by its EmployeeID.
func (h *adminCompanyHandler) GetAdminCompanyByEmployeeID(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	adminCompany, err := h.adminCompanyService.GetAdminCompanyByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if adminCompany == nil {
		helper.SendError(c, http.StatusNotFound, "Admin company not found for this employee ID")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Admin company fetched successfully", adminCompany)
}

// ChangeAdminPassword handles changing the password for the logged-in admin.
func (h *adminCompanyHandler) ChangeAdminPassword(c *gin.Context) {
	adminID, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Admin ID not found in token claims.")
		return
	}

	admID, ok := adminID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid admin ID type in token claims.")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	if err := h.adminCompanyService.ChangeAdminPassword(int(admID), req.OldPassword, req.NewPassword); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Password changed successfully.", nil)
}

// CheckAndNotifySubscriptions checks subscription statuses and sends notifications.
func (h *adminCompanyHandler) CheckAndNotifySubscriptions(c *gin.Context) {
	if err := h.adminCompanyService.CheckAndNotifySubscriptions(); err != nil {
		log.Printf("Error checking and notifying subscriptions: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to process subscription checks and notifications")
		return
	}
	helper.SendSuccess(c, http.StatusOK, "Subscription check and notifications processed.", nil)
}

// GetDashboardSummary handles fetching summary data for the admin dashboard.
func (h *adminCompanyHandler) GetDashboardSummary(hub *websocket.Hub, c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}

	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	summaryData, err := h.adminCompanyService.GetDashboardSummaryData(int(compID))
	if err != nil {
		log.Printf("Error getting dashboard summary data: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve dashboard summary.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Dashboard summary fetched successfully.", summaryData)

	// Send update to WebSocket clients
	hub.SendDashboardUpdate(int(compID), summaryData)
}

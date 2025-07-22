package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go-face-auth/helper"

	"go-face-auth/services"
	"go-face-auth/websocket"

	"log"


	"github.com/gin-gonic/gin"
)

// Activity represents a single recent activity for the dashboard.
type Activity struct {
	Type        string    `json:"type"` // e.g., "attendance", "leave_request"
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// --- Company Handlers ---

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateCompany(c *gin.Context) {
	var req services.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	company, err := services.CreateCompany(req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create company.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Company created successfully.", company)
}

func GetCompanyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	company, err := services.GetCompanyByID(id)
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

// GetDashboardSummary handles fetching summary data for the admin dashboard.
func GetDashboardSummary(hub *websocket.Hub, c *gin.Context) {
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

	summaryData, err := services.GetDashboardSummaryData(int(compID))
	if err != nil {
		log.Printf("Error getting dashboard summary data: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve dashboard summary.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Dashboard summary fetched successfully.", summaryData)

	// Send update to WebSocket clients
	hub.SendDashboardUpdate(int(compID), summaryData)
}
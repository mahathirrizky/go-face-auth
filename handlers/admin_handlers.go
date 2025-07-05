package handlers

import (
	"net/http"
	"strconv"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"

	"github.com/gin-gonic/gin"
	"log"

)

// --- Company Handlers ---

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	company := &models.CompaniesTable{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := repository.CreateCompany(company); err != nil {
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

	company, err := repository.GetCompanyByID(id)
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
func GetDashboardSummary(c *gin.Context) {
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

	// Fetch total employees
	totalEmployees, err := repository.GetTotalEmployeesByCompanyID(int(compID))
	if err != nil {
		log.Printf("Error getting total employees for company %d: %v", int(compID), err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve total employees.")
		return
	}

	// Fetch today's attendance (present, absent, leave)
	// You'll need to implement these repository functions
	presentToday, err := repository.GetPresentEmployeesCountToday(int(compID))
	if err != nil {
		log.Printf("Error getting present employees today for company %d: %v", int(compID), err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve today's attendance.")
		return
	}

	absentToday, err := repository.GetAbsentEmployeesCountToday(int(compID))
	if err != nil {
		log.Printf("Error getting absent employees today for company %d: %v", int(compID), err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve today's attendance.")
		return
	}

	onLeaveToday, err := repository.GetOnLeaveEmployeesCountToday(int(compID))
	if err != nil {
		log.Printf("Error getting on leave employees today for company %d: %v", int(compID), err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve today's attendance.")
		return
	}

	summary := gin.H{
		"total_employees": totalEmployees,
		"present_today":   presentToday,
		"absent_today":    absentToday,
		"on_leave_today":  onLeaveToday,
	}

	helper.SendSuccess(c, http.StatusOK, "Dashboard summary fetched successfully.", summary)
}



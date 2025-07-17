package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

type CreateLeaveRequestPayload struct {
	Type      string    `json:"type" binding:"required,oneof=cuti sakit"`
	StartDate string    `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string    `json:"end_date" binding:"required,datetime=2006-01-02"`
	Reason    string    `json:"reason" binding:"required,min=10"`
}

// Employee Handlers

func ApplyLeave(c *gin.Context) {
	employeeID, exists := c.Get("employeeID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token")
		return
	}
	empIDFloat, ok := employeeID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := uint(empIDFloat)

	var req CreateLeaveRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD.")
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD.")
		return
	}

	if startDate.After(endDate) {
		helper.SendError(c, http.StatusBadRequest, "Start date cannot be after end date.")
		return
	}

	// Ensure the employee exists
	_, err = repository.GetEmployeeByID(int(empID))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	leaveRequest := &models.LeaveRequest{
		EmployeeID: empID,
		Type:       req.Type,
		StartDate:  startDate,
		EndDate:    endDate,
		Reason:     req.Reason,
		Status:     "pending", // Default status
	}

	if err := repository.CreateLeaveRequest(leaveRequest); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to submit leave request.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Leave request submitted successfully.", leaveRequest)
}

func GetMyLeaveRequests(c *gin.Context) {
	employeeID, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token")
		return
	}
	empIDFloat, ok := employeeID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := uint(empIDFloat)

	var startDate *time.Time
	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD.")
			return
		}
		startDate = &parsedDate
	}

	var endDate *time.Time
	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD.")
			return
		}
		endDate = &parsedDate
	}

	leaveRequests, err := repository.GetLeaveRequestsByEmployeeID(empID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve leave requests.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Leave requests retrieved successfully.", leaveRequests)
}

// Admin Handlers

func GetAllCompanyLeaveRequests(c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}
	compID := int(compIDFloat)

	// Get query params for pagination and filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")

	leaveRequests, totalRecords, err := repository.GetCompanyLeaveRequestsPaginated(compID, status, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve leave requests.")
		return
	}

	paginatedData := gin.H{
		"items":         leaveRequests,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Leave requests retrieved successfully.", paginatedData)
}

type ReviewLeaveRequestPayload struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

func ReviewLeaveRequest(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
	leaveRequestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid leave request ID.")
		return
	}

	adminID, exists := c.Get("id") // Assuming adminID is set in JWT for admin users
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Admin ID not found in token.")
		return
	}
	adminIDUint, ok := adminID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid admin ID type in token claims.")
		return
	}
	adminIDVal := uint(adminIDUint)

	var req ReviewLeaveRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
		return
	}

	leaveRequest, err := repository.GetLeaveRequestByID(uint(leaveRequestID))
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Leave request not found.")
		return
	}

	// Ensure the admin reviewing is from the same company as the employee
	// First, get the employee's company ID
	employee, err := repository.GetEmployeeByID(int(leaveRequest.EmployeeID))
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusInternalServerError, "Could not find employee for leave request.")
		return
	}

	// Then, get the admin's company ID
	adminCompany, err := repository.GetAdminCompanyByID(int(adminIDVal))
	if err != nil || adminCompany == nil || adminCompany.CompanyID != employee.CompanyID {
		helper.SendError(c, http.StatusForbidden, "You are not authorized to review this leave request.")
		return
	}

	leaveRequest.Status = req.Status
	leaveRequest.ReviewedBy = &adminIDVal
	now := time.Now()
	leaveRequest.ReviewedAt = &now

	if err := repository.UpdateLeaveRequest(leaveRequest); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update leave request status.")
		return
	}

	// Get company ID from admin's token
	companyID, exists := c.Get("companyID")
	if !exists {
		return // Should not happen if AuthMiddleware is used
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		return // Should not happen
	}
	compID := int(compIDFloat)

	// Trigger dashboard update for the company
	go func() {
		summary, err := GetDashboardSummaryData(compID)
		if err != nil {
			log.Printf("Error fetching dashboard summary for WebSocket update after leave review: %v", err)
			return
		}
		hub.SendDashboardUpdate(compID, summary)
	}()

	helper.SendSuccess(c, http.StatusOK, "Leave request status updated successfully.", leaveRequest)
}
}
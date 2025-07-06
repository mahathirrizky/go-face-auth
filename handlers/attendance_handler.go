package handlers

import (
	"log"
	"net/http"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// AttendanceRequest represents the request body for attendance.
type AttendanceRequest struct {
	EmployeeID int `json:"employee_id" binding:"required"`
}

// HandleAttendance handles check-in and check-out processes.
func HandleAttendance(hub *websocket.Hub, c *gin.Context) {
	var req AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	// Check if employee exists
	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get latest attendance record for this employee
	latestAttendance, err := repository.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve attendance record.")
		return
	}

	now := time.Now()
	var message string

	if latestAttendance != nil && latestAttendance.CheckOutTime == nil {
		// Employee is currently checked in, so this is a check-out
		latestAttendance.CheckOutTime = &now
		latestAttendance.Status = "present" // Assuming successful check-out
		err = repository.UpdateAttendance(latestAttendance)
		message = "Check-out successful!"
	} else {
		// Employee is not checked in, so this is a check-in
		newAttendance := &models.AttendancesTable{
			EmployeeID:  req.EmployeeID,
			CheckInTime: now,
			Status:      "present", // Default status, can be refined later (e.g., 'late')
		}
		err = repository.CreateAttendance(newAttendance)
		message = "Check-in successful!"
	}

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record attendance.")
		return	
	}

	// Get company ID from JWT claims
	companyID, exists := c.Get("companyID")
	if !exists {
		// This should ideally not happen if AuthMiddleware is used
		return
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		// This should ideally not happen if AuthMiddleware is used
		return
	}
	compID := int(compIDFloat)

	// Fetch updated dashboard summary and send via WebSocket
	go func() {
		summary, err := GetDashboardSummaryData(compID)
		if err != nil {
			log.Printf("Error fetching dashboard summary for WebSocket update: %v", err)
			return
		}
		hub.SendDashboardUpdate(compID, summary)
	}()

	helper.SendSuccess(c, http.StatusOK, message, gin.H{
		"employee_id": employee.ID,
		"employee_name": employee.Name,
		"timestamp": now,
	})
}

// GetAttendances retrieves all attendance records for the company.
func GetAttendances(c *gin.Context) {
	// Get company ID from JWT claims
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token")
		return
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}
	compID := int(compIDFloat)

	attendances, err := repository.GetAttendancesByCompanyID(compID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve attendances.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Attendances retrieved successfully.", attendances)
}

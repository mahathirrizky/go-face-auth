package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// GetEmployeeAttendanceHistory retrieves attendance records for a specific employee.
func GetEmployeeAttendanceHistory(c *gin.Context) {
	employeeID := c.Param("employeeID")
	parsedEmployeeID, err := strconv.Atoi(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD.")
			return
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD.")
			return
		}
		// Set end date to end of day for inclusive range
		endDateVal := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &endDateVal
	}

	attendances, err := repository.GetEmployeeAttendances(parsedEmployeeID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee attendance history.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee attendance history retrieved successfully.", attendances)
}

// ExportEmployeeAttendanceToExcel exports attendance records for a specific employee to an Excel file.
func ExportEmployeeAttendanceToExcel(c *gin.Context) {
	employeeID := c.Param("employeeID")
	parsedEmployeeID, err := strconv.Atoi(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD.")
			return
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD.")
			return
		}
		endDateVal := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &endDateVal
	}

	attendances, err := repository.GetEmployeeAttendances(parsedEmployeeID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee attendance for export.")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	// Set headers
	f.SetCellValue("Sheet1", "A1", "Employee Name")
	f.SetCellValue("Sheet1", "B1", "Check In Time")
	f.SetCellValue("Sheet1", "C1", "Check Out Time")
	f.SetCellValue("Sheet1", "D1", "Status")

	// Populate data
	for i, att := range attendances {
		row := i + 2 // Start from row 2 after headers
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), att.Employee.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), att.CheckInTime.Format("2006-01-02 15:04:05"))
		checkOutTime := "N/A"
		if att.CheckOutTime != nil {
			checkOutTime = att.CheckOutTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), checkOutTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), att.Status)
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employee_attendance.xlsx")

	// Write the Excel file to the response writer
	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// ExportAllAttendancesToExcel exports all attendance records for the company to an Excel file.
func ExportAllAttendancesToExcel(c *gin.Context) {
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

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD.")
			return
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD.")
			return
		}
		endDateVal := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &endDateVal
	}

	attendances, err := repository.GetCompanyAttendancesFiltered(compID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve all company attendances for export.")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	// Set headers
	f.SetCellValue("Sheet1", "A1", "Employee Name")
	f.SetCellValue("Sheet1", "B1", "Check In Time")
	f.SetCellValue("Sheet1", "C1", "Check Out Time")
	f.SetCellValue("Sheet1", "D1", "Status")

	// Populate data
	for i, att := range attendances {
		row := i + 2 // Start from row 2 after headers
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), att.Employee.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), att.CheckInTime.Format("2006-01-02 15:04:05"))
		checkOutTime := "N/A"
		if att.CheckOutTime != nil {
			checkOutTime = att.CheckOutTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), checkOutTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), att.Status)
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=all_company_attendance.xlsx")

	// Write the Excel file to the response writer
	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

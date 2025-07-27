package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-face-auth/helper"
	"go-face-auth/services"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
)

// AttendanceHandler defines the interface for attendance related handlers.
type AttendanceHandler interface {
	HandleAttendance(hub *websocket.Hub, c *gin.Context)
	HandleOvertimeCheckIn(hub *websocket.Hub, c *gin.Context)
	HandleOvertimeCheckOut(hub *websocket.Hub, c *gin.Context)
	GetAttendances(c *gin.Context)
	GetEmployeeAttendanceHistory(c *gin.Context)
	ExportEmployeeAttendanceToExcel(c *gin.Context)
	ExportAllAttendancesToExcel(c *gin.Context)
	GetUnaccountedEmployees(c *gin.Context)
	ExportUnaccountedToExcel(c *gin.Context)
	ExportOvertimeToExcel(c *gin.Context)
	GetOvertimeAttendances(c *gin.Context)
	CorrectAttendance(c *gin.Context)
}

// attendanceHandler is the concrete implementation of AttendanceHandler.
type attendanceHandler struct {
	attendanceService services.AttendanceService
	adminCompanyService services.AdminCompanyService
}

// NewAttendanceHandler creates a new instance of AttendanceHandler.
func NewAttendanceHandler(attendanceService services.AttendanceService, adminCompanyService services.AdminCompanyService) AttendanceHandler {
	return &attendanceHandler{
		attendanceService: attendanceService,
		adminCompanyService:      adminCompanyService,
	}
}

// HandleAttendance handles regular check-in and check-out processes.
func (h *attendanceHandler) HandleAttendance(hub *websocket.Hub, c *gin.Context) {
	var req services.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	message, employee, now, err := h.attendanceService.HandleAttendance(req)
	if err != nil {
		// Check for specific error messages from the service
		if err.Error() == "face not recognized" {
			helper.SendError(c, http.StatusConflict, err.Error())
		} else if err.Error() == "no registered face images for this employee" || err.Error() == "employee not found" {
			helper.SendError(c, http.StatusNotFound, err.Error())
		} else if err.Error() == "you are not within a valid attendance location" || err.Error() == "employee does not have a shift assigned" || err.Error() == "cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead" {
			helper.SendError(c, http.StatusBadRequest, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	companyID, exists := c.Get("companyID")
	if !exists {
		return
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		return
	}
	compID := int(compIDFloat)

	go func() {
		summary, err := h.adminCompanyService.GetDashboardSummaryData(compID)
		if err != nil {
			log.Printf("Error fetching dashboard summary for WebSocket update: %v", err)
			return
		}
		hub.SendDashboardUpdate(compID, summary)
	}()

	helper.SendSuccess(c, http.StatusOK, message, gin.H{
		"employee_id":   employee.ID,
		"employee_name": employee.Name,
		"timestamp":     now,
	})
}

// HandleOvertimeCheckIn handles overtime check-in process.
func (h *attendanceHandler) HandleOvertimeCheckIn(hub *websocket.Hub, c *gin.Context) {
	var req services.OvertimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	employee,now, err := h.attendanceService.HandleOvertimeCheckIn(req)
	if err != nil {
		// Check for specific error messages from the service
		if err.Error() == "face not recognized" {
			helper.SendError(c, http.StatusConflict, err.Error())
		} else if err.Error() == "no registered face images for this employee" || err.Error() == "employee not found" {
			helper.SendError(c, http.StatusNotFound, err.Error())
		} else if err.Error() == "you are not within a valid attendance location" || err.Error() == "employee does not have a shift assigned" || err.Error() == "cannot check-in for overtime during regular shift hours" || err.Error() == "employee is already checked in for overtime" {
			helper.SendError(c, http.StatusBadRequest, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Overtime check-in successful!", gin.H{
		"employee_id":   employee.ID,
		"employee_name": employee.Name,
		"timestamp":     now,
		"status":        "overtime_in",
	})
}

// HandleOvertimeCheckOut handles overtime check-out process.
func (h *attendanceHandler) HandleOvertimeCheckOut(hub *websocket.Hub, c *gin.Context) {
	var req services.OvertimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	employee, CheckInTime, now, OvertimeMinutes, err := h.attendanceService.HandleOvertimeCheckOut(req)
	if err != nil {
		// Check for specific error messages from the service
		if err.Error() == "face not recognized" {
			helper.SendError(c, http.StatusConflict, err.Error())
		} else if err.Error() == "no registered face images for this employee" || err.Error() == "employee not found" {
			helper.SendError(c, http.StatusNotFound, err.Error())
		} else if err.Error() == "employee is not currently checked in for overtime" {
			helper.SendError(c, http.StatusBadRequest, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Overtime check-out successful!", gin.H{
		"employee_id":      employee.ID,
		"employee_name":    employee.Name,
		"check_in_time":    CheckInTime,
		"check_out_time":   now,
		"overtime_minutes": OvertimeMinutes,
		"status":           "overtime_out",
	})
}

// GetAttendances retrieves all attendance records for the company.
func (h *attendanceHandler) GetAttendances(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

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
		endDate = &parsed
	}

	attendances, totalRecords, err := h.attendanceService.GetAttendancesPaginated(compID, startDate, endDate, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve attendances.")
		return
	}

	paginatedData := gin.H{
		"items":         attendances,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Attendances retrieved successfully.", paginatedData)
}

// GetEmployeeAttendanceHistory retrieves attendance records for a specific employee.
func (h *attendanceHandler) GetEmployeeAttendanceHistory(c *gin.Context) {
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

	attendances, err := h.attendanceService.GetEmployeeAttendances(parsedEmployeeID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee attendance history.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee attendance history retrieved successfully.", attendances)
}

// ExportEmployeeAttendanceToExcel exports attendance records for a specific employee to an Excel file.
func (h *attendanceHandler) ExportEmployeeAttendanceToExcel(c *gin.Context) {
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

	file, fileName, err := h.attendanceService.ExportEmployeeAttendanceToExcel(parsedEmployeeID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// Write the Excel file to the response writer
	if err := file.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// ExportAllAttendancesToExcel exports all attendance records for the company to an Excel file.
func (h *attendanceHandler) ExportAllAttendancesToExcel(c *gin.Context) {
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

	file, fileName, err := h.attendanceService.ExportAllAttendancesToExcel(compID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// Write the Excel file to the response writer
	if err := file.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// GetUnaccountedEmployees handles fetching employees who are not present and not on leave/sick.
func (h *attendanceHandler) GetUnaccountedEmployees(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

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
		endDate = &parsed
	}

	unaccountedEmployees, totalRecords, err := h.attendanceService.GetUnaccountedEmployeesPaginated(compID, startDate, endDate, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve unaccounted employees.")
		return
	}

	paginatedData := gin.H{
		"items":         unaccountedEmployees,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Unaccounted employees retrieved successfully.", paginatedData)
}

// ExportUnaccountedToExcel exports unaccounted employee records to an Excel file.
func (h *attendanceHandler) ExportUnaccountedToExcel(c *gin.Context) {
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
		endDate = &parsed
	}

	search := c.Query("search")

	file, fileName, err := h.attendanceService.ExportUnaccountedToExcel(compID, startDate, endDate, search)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := file.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// ExportOvertimeToExcel exports overtime attendance records to an Excel file.
func (h *attendanceHandler) ExportOvertimeToExcel(c *gin.Context) {
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
		endDate = &parsed
	}

	search := c.Query("search")

	file, fileName, err := h.attendanceService.ExportOvertimeToExcel(compID, startDate, endDate, search)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := file.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// GetOvertimeAttendances retrieves all overtime attendance records for the company.
func (h *attendanceHandler) GetOvertimeAttendances(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

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
		endDate = &parsed
	}

	overtimeAttendances, totalRecords, err := h.attendanceService.GetOvertimeAttendancesPaginated(compID, startDate, endDate, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve overtime attendances.")
		return
	}

	paginatedData := gin.H{
		"items":         overtimeAttendances,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Overtime attendances retrieved successfully.", paginatedData)
}

// CorrectAttendance handles manual attendance correction by an admin.
func (h *attendanceHandler) CorrectAttendance(c *gin.Context) {
	adminID, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Admin ID not found in token.")
		return
	}
	adminIDUint, ok := adminID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid admin ID type in token claims.")
		return
	}

	var req services.CorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	attendance, err := h.attendanceService.CorrectAttendance(uint(adminIDUint), req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Attendance corrected successfully.", attendance)
}

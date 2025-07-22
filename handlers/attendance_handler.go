package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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
	EmployeeID int     `json:"employee_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	ImageData  string  `json:"image_data" binding:"required"`
}

// OvertimeAttendanceRequest represents the request body for overtime attendance.
type OvertimeAttendanceRequest struct {
	EmployeeID int     `json:"employee_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	ImageData  string  `json:"image_data" binding:"required"`
}

// PythonRecognitionRequest to Python server
type PythonRecognitionRequest struct {
	ClientImageData string `json:"client_image_data"` // Base64 encoded image from client
	DBImagePath     string `json:"db_image_path"`     // Path to the image file on the Python server's side
}

// HandleAttendance handles regular check-in and check-out processes.
func HandleAttendance(hub *websocket.Hub, c *gin.Context) {
	var req AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		helper.SendError(c, http.StatusInternalServerError, "Could not retrieve employee face image.")
		return
	}
	if len(faceImages) == 0 {
		helper.SendError(c, http.StatusNotFound, "No registered face images for this employee.")
		return
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Face recognition service is unavailable.")
		return
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		helper.SendError(c, http.StatusConflict, "Face not recognized.")
		return
	}
	// --- End of Face Recognition Logic ---

	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information.")
		return
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		helper.SendError(c, http.StatusInternalServerError, "Invalid company timezone configuration.")
		return
	}
	

	// Get all valid attendance locations for the company
	companyLocations, err := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil || len(companyLocations) == 0 {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company attendance locations or no locations configured.")
		return
	}

	// Validate employee's current location against company's valid attendance locations
	isWithinValidLocation := false
	for _, loc := range companyLocations {
		distance := helper.HaversineDistance(req.Latitude, req.Longitude, loc.Latitude, loc.Longitude)
		if distance <= float64(loc.Radius) {
			isWithinValidLocation = true
			break
		}
	}

	if !isWithinValidLocation {
		helper.SendError(c, http.StatusBadRequest, "You are not within a valid attendance location.")
		return
	}

	// Get employee's shift
	if employee.ShiftID == nil {
		helper.SendError(c, http.StatusBadRequest, "Employee does not have a shift assigned.")
		return
	}
	shift := employee.Shift // Shift is preloaded by GetEmployeeByID

	now := time.Now().In(companyLocation) // Get current time in company's timezone
	var message string
	
	latestAttendance, err := repository.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve attendance record.")
		return
	}

	if latestAttendance != nil && latestAttendance.CheckOutTime == nil {
		// Regular Check-out
		latestAttendance.CheckOutTime = &now
		latestAttendance.Status = "present" 
		err = repository.UpdateAttendance(latestAttendance)
		message = "Check-out successful!"
	} else {
		// Regular Check-in
		// Check if current time is within regular shift
		isWithinShift, err := helper.IsTimeWithinShift(now, shift.StartTime, shift.EndTime, shift.GracePeriodMinutes, companyLocation)
		if err != nil {
			log.Printf("Error checking time within shift: %v", err)
			helper.SendError(c, http.StatusInternalServerError, "Failed to validate shift time.")
			return
		}

		if !isWithinShift {
			helper.SendError(c, http.StatusBadRequest, "Cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead.")
			return
		}

		// Determine status (on time or late)
		shiftStartToday, _ := helper.ParseTime(now, shift.StartTime, companyLocation)
		if now.After(shiftStartToday.Add(time.Duration(shift.GracePeriodMinutes) * time.Minute)) {
			status = "late"
		} else {
			status = "on_time"
		}

		newAttendance := &models.AttendancesTable{
			EmployeeID:  req.EmployeeID,
			CheckInTime: now,
			Status:      status,
		}
		err = repository.CreateAttendance(newAttendance)
		if err != nil {	
			helper.SendError(c, http.StatusInternalServerError, "Failed to record attendance.")	
			return	
		}	
		message = "Check-in successful!"
	}

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record attendance.")
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
		summary, err := GetDashboardSummaryData(compID)
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
func HandleOvertimeCheckIn(hub *websocket.Hub, c *gin.Context) {
	var req OvertimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		helper.SendError(c, http.StatusInternalServerError, "Could not retrieve employee face image.")
		return
	}
	if len(faceImages) == 0 {
		helper.SendError(c, http.StatusNotFound, "No registered face images for this employee.")
		return
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Face recognition service is unavailable.")
		return
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		helper.SendError(c, http.StatusConflict, "Face not recognized.")
		return
	}
	// --- End of Face Recognition Logic ---

	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information.")
		return
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		helper.SendError(c, http.StatusInternalServerError, "Invalid company timezone configuration.")
		return
	}

	// Get all valid attendance locations for the company
	companyLocations, err := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil || len(companyLocations) == 0 {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company attendance locations or no locations configured.")
		return
	}

	// Validate employee's current location against company's valid attendance locations
	isWithinValidLocation := false
	for _, loc := range companyLocations {
		distance := helper.HaversineDistance(req.Latitude, req.Longitude, loc.Latitude, loc.Longitude)
		if distance <= float64(loc.Radius) {
			isWithinValidLocation = true
			break
		}
	}

	if !isWithinValidLocation {
		helper.SendError(c, http.StatusBadRequest, "You are not within a valid attendance location.")
		return
	}

	// Get employee's shift
	if employee.ShiftID == nil {
		helper.SendError(c, http.StatusBadRequest, "Employee does not have a shift assigned.")
		return
	}
	shift := employee.Shift

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Validate: Cannot check-in for overtime if within regular shift hours
	isWithinShift, err := helper.IsTimeWithinShift(now, shift.StartTime, shift.EndTime, shift.GracePeriodMinutes, companyLocation)
	if err != nil {
		log.Printf("Error checking time within shift for overtime check-in: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to validate shift time.")
		return
	}
	if isWithinShift {
		helper.SendError(c, http.StatusBadRequest, "Cannot check-in for overtime during regular shift hours.")
		return
	}

	// Check if employee is already checked in for overtime
	latestOvertimeAttendance, err := repository.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve latest overtime record.")
		return
	}
	if latestOvertimeAttendance != nil && latestOvertimeAttendance.CheckOutTime == nil && latestOvertimeAttendance.Status == "overtime_in" {
		helper.SendError(c, http.StatusBadRequest, "Employee is already checked in for overtime.")
		return
	}

	// Create new overtime check-in record
	newOvertimeAttendance := &models.AttendancesTable{
		EmployeeID:  req.EmployeeID,
		CheckInTime: now,
		Status:      "overtime_in", // Specific status for overtime check-in
	}
	err = repository.CreateAttendance(newOvertimeAttendance)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record overtime check-in.")
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
func HandleOvertimeCheckOut(hub *websocket.Hub, c *gin.Context) {
	var req OvertimeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		helper.SendError(c, http.StatusInternalServerError, "Could not retrieve employee face image.")
		return
	}
	if len(faceImages) == 0 {
		helper.SendError(c, http.StatusNotFound, "No registered face images for this employee.")
		return
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Face recognition service is unavailable.")
		return
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		helper.SendError(c, http.StatusConflict, "Face not recognized.")
		return
	}
	// --- End of Face Recognition Logic ---

	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information.")
		return
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		helper.SendError(c, http.StatusInternalServerError, "Invalid company timezone configuration.")
		return
	}

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Find the latest "overtime_in" record that is not checked out
	latestOvertimeAttendance, err := repository.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve latest overtime record.")
		return
	}
	if latestOvertimeAttendance == nil || latestOvertimeAttendance.CheckOutTime != nil || latestOvertimeAttendance.Status != "overtime_in" {
		helper.SendError(c, http.StatusBadRequest, "Employee is not currently checked in for overtime.")
		return
	}

	// Calculate overtime duration
	overtimeDuration := now.Sub(latestOvertimeAttendance.CheckInTime)
	overtimeMinutes := int(overtimeDuration.Minutes())

	latestOvertimeAttendance.CheckOutTime = &now
	latestOvertimeAttendance.OvertimeMinutes = overtimeMinutes
	latestOvertimeAttendance.Status = "overtime_out" // Specific status for overtime check-out

	err = repository.UpdateAttendance(latestOvertimeAttendance)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record overtime check-out.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Overtime check-out successful!", gin.H{
		"employee_id":     employee.ID,
		"employee_name":   employee.Name,
		"check_in_time":   latestOvertimeAttendance.CheckInTime,
		"check_out_time":  now,
		"overtime_minutes": overtimeMinutes,
		"status":          "overtime_out",
	})
}

// GetAttendances retrieves all attendance records for the company.
func GetAttendances(c *gin.Context) {
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

	attendances, totalRecords, err := repository.GetAttendancesPaginated(compID, startDate, endDate, search, page, pageSize)
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

	// Define and set the sheet name
	sheetName := "Employee Attendance"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Check In Time")
	f.SetCellValue(sheetName, "C1", "Check Out Time")
	f.SetCellValue(sheetName, "D1", "Status")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		log.Printf("Error creating style: %v", err)
	} else {
		f.SetCellStyle(sheetName, "A1", "D1", style)
	}

	// Populate data
	for i, att := range attendances {
		row := i + 2 // Start from row 2 after headers
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), att.Employee.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), att.CheckInTime.Format("2006-01-02 15:04:05"))
		checkOutTime := "N/A"
		if att.CheckOutTime != nil {
			checkOutTime = att.CheckOutTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), checkOutTime)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), att.Status)
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	fileName := "employee_attendance.xlsx"
	if len(attendances) > 0 {
		employeeName := attendances[0].Employee.Name
		dateRange := ""
		if startDate != nil && endDate != nil {
			dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		} else if startDate != nil {
			dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
		} else if endDate != nil {
			dateRange = fmt.Sprintf("_until_%s", endDate.Format("2006-01-02"))
		}
		fileName = fmt.Sprintf("%s_attendance%s.xlsx", employeeName, dateRange)
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

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

	attendances, err := repository.GetCompanyAttendancesFiltered(compID, startDate, endDate, "all")
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

	// Define and set the sheet name
	sheetName := "All Attendances"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Check In Time")
	f.SetCellValue(sheetName, "C1", "Check Out Time")
	f.SetCellValue(sheetName, "D1", "Status")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		log.Printf("Error creating style: %v", err)
	} else {
		f.SetCellStyle(sheetName, "A1", "D1", style)
	}

	// Populate data
	for i, att := range attendances {
		row := i + 2 // Start from row 2 after headers
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), att.Employee.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), att.CheckInTime.Format("2006-01-02 15:04:05"))
		checkOutTime := "N/A"
		if att.CheckOutTime != nil {
			checkOutTime = att.CheckOutTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), checkOutTime)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), att.Status)
	}

	// Set response headers for Excel file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	fileName := "all_company_attendance.xlsx"
	dateRange := ""
	if startDate != nil && endDate != nil {
		dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
	} else if endDate != nil {
		dateRange = fmt.Sprintf("_until_%s", endDate.Format("2006-01-02"))
	}
	fileName = fmt.Sprintf("all_company_attendance%s.xlsx", dateRange)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// Write the Excel file to the response writer
	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// GetUnaccountedEmployees handles fetching employees who are not present and not on leave/sick.
func GetUnaccountedEmployees(c *gin.Context) {
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

	unaccountedEmployees, totalRecords, err := repository.GetUnaccountedEmployeesPaginated(compID, startDate, endDate, search, page, pageSize)
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
func ExportUnaccountedToExcel(c *gin.Context) {
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

	unaccountedEmployees, err := repository.GetUnaccountedEmployeesFiltered(compID, startDate, endDate, search)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve unaccounted employees for export.")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	sheetName := "Unaccounted Employees"
	f.SetSheetName("Sheet1", sheetName)

	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Email")
	f.SetCellValue(sheetName, "C1", "Position")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		log.Printf("Error creating style: %v", err)
	} else {
		f.SetCellStyle(sheetName, "A1", "C1", style)
	}

	for i, emp := range unaccountedEmployees {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), emp.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), emp.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), emp.Position)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	fileName := "unaccounted_employees.xlsx"
	dateRange := ""
	if startDate != nil && endDate != nil {
		dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
	} else if endDate != nil {
		dateRange = fmt.Sprintf("_until_%s", endDate.Format("2006-01-02"))
	}
	fileName = fmt.Sprintf("unaccounted_employees%s.xlsx", dateRange)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

// ExportOvertimeToExcel exports overtime attendance records to an Excel file.
func ExportOvertimeToExcel(c *gin.Context) {
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

	overtimeAttendances, err := repository.GetOvertimeAttendancesFiltered(compID, startDate, endDate, search)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve overtime attendances for export.")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	sheetName := "Overtime Attendances"
	f.SetSheetName("Sheet1", sheetName)

	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Check In Time")
	f.SetCellValue(sheetName, "C1", "Check Out Time")
	f.SetCellValue(sheetName, "D1", "Overtime Minutes")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		log.Printf("Error creating style: %v", err)
	} else {
		f.SetCellStyle(sheetName, "A1", "D1", style)
	}

	for i, att := range overtimeAttendances {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), att.Employee.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), att.CheckInTime.Format("2006-01-02 15:04:05"))
		checkOutTime := "N/A"
		if att.CheckOutTime != nil {
			checkOutTime = att.CheckOutTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), checkOutTime)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), att.OvertimeMinutes)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	fileName := "overtime_attendances.xlsx"
	dateRange := ""
	if startDate != nil && endDate != nil {
		dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
	} else if endDate != nil {
		dateRange = fmt.Sprintf("_until_%s", endDate.Format("2006-01-02"))
	}
	fileName = fmt.Sprintf("overtime_attendances%s.xlsx", dateRange)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}



// GetOvertimeAttendances retrieves all overtime attendance records for the company.
func GetOvertimeAttendances(c *gin.Context) {
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

	overtimeAttendances, totalRecords, err := repository.GetOvertimeAttendancesPaginated(compID, startDate, endDate, search, page, pageSize)
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

// sendToPythonServer connects to the Python TCP server, sends the payload, and returns the response.
func sendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error) {
	pythonServerAddr := "127.0.0.1:5000" // Python server address
	conn, err := net.Dial("tcp", pythonServerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Send payload to Python server with a newline delimiter
	_, err = conn.Write(append(payloadBytes, '\n'))
	if err != nil {
		return nil, err
	}

	// Read response from Python server
	decoder := json.NewDecoder(conn)
	var pythonResponse map[string]interface{}
	if err := decoder.Decode(&pythonResponse); err != nil {
		return nil, err
	}

	return pythonResponse, nil
}

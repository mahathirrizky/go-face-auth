package services

import (
	"encoding/json"
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"net"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

// AttendanceRequest represents the request body for attendance.
type AttendanceRequest struct {
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

// sendToPythonServer connects to the Python TCP server, sends the payload, and returns the response.
func sendToPythonServer(payload PythonRecognitionRequest) (map[string]interface{}, error) {
	pythonServerAddr := os.Getenv("PYTHON_SERVER_ADDRESS")
	if pythonServerAddr == "" {
		pythonServerAddr = "127.0.0.1:5000" // Default to localhost if not set
	}
	conn, err := net.Dial("tcp", pythonServerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Python server: %w", err)
	}
	defer conn.Close()

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send payload to Python server with a newline delimiter
	_, err = conn.Write(append(payloadBytes, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to send payload to Python server: %w", err)
	}

	// Read response from Python server
	decoder := json.NewDecoder(conn)
	var pythonResponse map[string]interface{}
	if err := decoder.Decode(&pythonResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response from Python server: %w", err)
	}

	return pythonResponse, nil
}

func  HandleAttendance(req AttendanceRequest) (string, *models.EmployeesTable, time.Time, error) {
	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return "", nil, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company information")
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return "", nil, time.Time{}, fmt.Errorf("invalid company timezone configuration")
	}

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Check if employee is on approved leave for today
	approvedLeave, err := repository.IsEmployeeOnApprovedLeave(employee.ID, now)
	if err != nil {
		log.Printf("Error checking leave status for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
		return "", nil, time.Time{}, fmt.Errorf("failed to check leave status")
	}

	if approvedLeave != nil {
		var leaveType string
		if approvedLeave.Type == "sakit" {
			leaveType = "sakit"
		} else {
			leaveType = "cuti"
		}
		return "", nil, time.Time{}, fmt.Errorf("anda sedang dalam pengajuan %s yang disetujui untuk hari ini", leaveType)
	}

	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		return "", nil, time.Time{}, fmt.Errorf("could not retrieve employee face image")
	}
	if len(faceImages) == 0 {
		return "", nil, time.Time{}, fmt.Errorf("no registered face images for this employee")
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return "", nil, time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return "", nil, time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	employee, err = repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return "", nil, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err = repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company information")
	}

	companyLocation, err = time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return "", nil, time.Time{}, fmt.Errorf("invalid company timezone configuration")
	}

	// Get all valid attendance locations for the company
	companyLocations, err := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil || len(companyLocations) == 0 {
		return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations or no locations configured")
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
		return "", nil, time.Time{}, fmt.Errorf("you are not within a valid attendance location")
	}

	// Get employee's shift
	if employee.ShiftID == nil {
		return "", nil, time.Time{}, fmt.Errorf("employee does not have a shift assigned")
	}
	shift := employee.Shift // Shift is preloaded by GetEmployeeByID

	now = time.Now().In(companyLocation) // Get current time in company's timezone
	var message string

	latestAttendance, err := repository.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to retrieve attendance record")
	}

	if latestAttendance != nil && latestAttendance.CheckOutTime == nil {
		// Regular Check-out
		latestAttendance.CheckOutTime = &now
		latestAttendance.Status = "present"
		err = repository.UpdateAttendance(latestAttendance)
		message = "Check-out successful!"
	} else {
		// Regular Check-in
		// Calculate earliest allowed check-in time (1.5 hours before shift start)
		shiftStartToday, err := helper.ParseTime(now, shift.StartTime, companyLocation)
		if err != nil {
			log.Printf("Error parsing shift start time for early check-in: %v", err)
			return "", nil, time.Time{}, fmt.Errorf("failed to validate shift time")
		}
		earliesCheckInTime := shiftStartToday.Add(-90 * time.Minute) // 90 minutes = 1.5 hours

		// Prevent check-in if too early
		if now.Before(earliesCheckInTime) {
			return "", nil, time.Time{}, fmt.Errorf("Anda tidak dapat absen lebih dari 1.5 jam sebelum jam shift Anda.")
		}

		// Check if current time is within regular shift (considering grace period for late check-in)
		isWithinShift, err := helper.IsTimeWithinShift(now, shift.StartTime, shift.EndTime, shift.GracePeriodMinutes, companyLocation)
		if err != nil {
			log.Printf("Error checking time within shift: %v", err)
			return "", nil, time.Time{}, fmt.Errorf("failed to validate shift time")
		}

		if !isWithinShift {
			return "", nil, time.Time{}, fmt.Errorf("cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead")
		}

		// Determine status (on time or late)
		// shiftStartToday, _ := helper.ParseTime(now, shift.StartTime, companyLocation) // Already parsed above
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
			return "", nil, time.Time{}, fmt.Errorf("failed to record attendance")
		}
		message = "Check-in successful!"
	}

	if err != nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to record attendance")
	}

	return message, employee, now, nil
}

// OvertimeAttendanceRequest represents the request body for overtime attendance.
type OvertimeAttendanceRequest struct {
	EmployeeID int     `json:"employee_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	ImageData  string  `json:"image_data" binding:"required"`
}

// HandleOvertimeCheckIn handles overtime check-in process.
func HandleOvertimeCheckIn(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, error) {
	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		return nil,time.Time{}, fmt.Errorf("could not retrieve employee face image")
	}
	if len(faceImages) == 0 {
		return nil,time.Time{}, fmt.Errorf("no registered face images for this employee")
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return nil,time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return nil,time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil,time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		return nil,time.Time{}, fmt.Errorf("failed to retrieve company information")
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return nil,time.Time{}, fmt.Errorf("invalid company timezone configuration")
	}

	// Get all valid attendance locations for the company
	companyLocations, err := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil || len(companyLocations) == 0 {
		return nil,time.Time{}, fmt.Errorf("failed to retrieve company attendance locations or no locations configured")
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
		return nil,time.Time{}, fmt.Errorf("you are not within a valid attendance location")
	}

	// Get employee's shift
	if employee.ShiftID == nil {
		return nil,time.Time{}, fmt.Errorf("employee does not have a shift assigned")
	}
	shift := employee.Shift

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Validate: Cannot check-in for overtime if within regular shift hours
	isWithinShift, err := helper.IsTimeWithinShift(now, shift.StartTime, shift.EndTime, shift.GracePeriodMinutes, companyLocation)
	if err != nil {
		log.Printf("Error checking time within shift for overtime check-in: %v", err)
		return nil,time.Time{}, fmt.Errorf("failed to validate shift time")
	}
	if isWithinShift {
		return nil,time.Time{}, fmt.Errorf("cannot check-in for overtime during regular shift hours")
	}

	// Check if employee has an open regular check-in
	latestRegularAttendance, err := repository.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil,time.Time{}, fmt.Errorf("failed to retrieve latest regular attendance record")
	}
	if latestRegularAttendance != nil && latestRegularAttendance.CheckOutTime == nil && latestRegularAttendance.Status != "overtime_in" && latestRegularAttendance.Status != "overtime_out" {
		return nil,time.Time{}, fmt.Errorf("Anda harus check-out dari shift reguler sebelum check-in lembur.")
	}

	// Check if employee is already checked in for overtime
	latestOvertimeAttendance, err := repository.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil,time.Time{}, fmt.Errorf("failed to retrieve latest overtime record")
	}
	if latestOvertimeAttendance != nil && latestOvertimeAttendance.CheckOutTime == nil && latestOvertimeAttendance.Status == "overtime_in" {
		return nil,time.Time{}, fmt.Errorf("employee is already checked in for overtime")
	}

	// Create new overtime check-in record
	newOvertimeAttendance := &models.AttendancesTable{
		EmployeeID:  req.EmployeeID,
		CheckInTime: now,
		Status:      "overtime_in", // Specific status for overtime check-in
	}
	err = repository.CreateAttendance(newOvertimeAttendance)
	if err != nil {
		return nil,time.Time{}, fmt.Errorf("failed to record overtime check-in")
	}

	return employee, now,nil
}

func HandleOvertimeCheckOut(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, int, time.Time, error) {
	// --- Face Recognition Logic ---
	faceImages, err := repository.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("could not retrieve employee face image")
	}
	if len(faceImages) == 0 {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("no registered face images for this employee")
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := sendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := repository.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("failed to retrieve company information")
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("invalid company timezone configuration")
	}

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Find the latest "overtime_in" record that is not checked out
	latestOvertimeAttendance, err := repository.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("failed to retrieve latest overtime record")
	}
	if latestOvertimeAttendance == nil || latestOvertimeAttendance.CheckOutTime != nil || latestOvertimeAttendance.Status != "overtime_in" {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("employee is not currently checked in for overtime")
	}

	// Calculate overtime duration
	overtimeDuration := now.Sub(latestOvertimeAttendance.CheckInTime)
	overtimeMinutes := int(overtimeDuration.Minutes())

	latestOvertimeAttendance.CheckOutTime = &now
	latestOvertimeAttendance.OvertimeMinutes = overtimeMinutes
	latestOvertimeAttendance.Status = "overtime_out" // Specific status for overtime check-out

	err = repository.UpdateAttendance(latestOvertimeAttendance)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("failed to record overtime check-out")
	}

	// Return employee, now, overtimeMinutes, and original CheckInTime
	return employee, now, overtimeMinutes, latestOvertimeAttendance.CheckInTime, nil
}

func GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error) {
	return repository.GetAttendancesPaginated(companyID, startDate, endDate, search, page, pageSize)
}

func ExportEmployeeAttendanceToExcel(employeeID int, startDate, endDate *time.Time) (*excelize.File, string, error) {
	attendances, err := repository.GetEmployeeAttendances(employeeID, startDate, endDate)
	if err != nil {
		return nil, "", fmt.Errorf("failed to retrieve employee attendance for export: %w", err)
	}

	f := excelize.NewFile()
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
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font:      &excelize.Font{Bold: true},
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

	return f, fileName, nil
}

func ExportAllAttendancesToExcel(companyID int, startDate, endDate *time.Time) (*excelize.File, string, error) {
	attendances, err := repository.GetCompanyAttendancesFiltered(companyID, startDate, endDate, "all")
	if err != nil {
		return nil, "", fmt.Errorf("failed to retrieve all company attendances for export: %w", err)
	}

	f := excelize.NewFile()
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
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font:      &excelize.Font{Bold: true},
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

	return f, fileName, nil
}

func ExportUnaccountedToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error) {
	unaccountedEmployees, err := repository.GetUnaccountedEmployeesFiltered(companyID, startDate, endDate, search)
	if err != nil {
		return nil, "", fmt.Errorf("failed to retrieve unaccounted employees for export: %w", err)
	}

	f := excelize.NewFile()
	sheetName := "Unaccounted Employees"
	f.SetSheetName("Sheet1", sheetName)

	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Email")
	f.SetCellValue(sheetName, "C1", "Position")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font:      &excelize.Font{Bold: true},
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

	fileName := "unaccounted_employees.xlsx"
	dateRange := ""
	if startDate != nil && endDate != nil {
		dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
	}
	fileName = fmt.Sprintf("unaccounted_employees%s.xlsx", dateRange)

	return f, fileName, nil
}

func ExportOvertimeToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error) {
	overtimeAttendances, err := repository.GetOvertimeAttendancesFiltered(companyID, startDate, endDate, search)
	if err != nil {
		return nil, "", fmt.Errorf("failed to retrieve overtime attendances for export: %w", err)
	}

	f := excelize.NewFile()
	sheetName := "Overtime Attendances"
	f.SetSheetName("Sheet1", sheetName)

	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Check In Time")
	f.SetCellValue(sheetName, "C1", "Check Out Time")
	f.SetCellValue(sheetName, "D1", "Overtime Minutes")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font:      &excelize.Font{Bold: true},
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

	return f, fileName, nil
}

func GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error) {
	return repository.GetOvertimeAttendancesPaginated(companyID, startDate, endDate, search, page, pageSize)
}
func GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
	return repository.GetUnaccountedEmployeesPaginated(companyID, startDate, endDate, search, page, pageSize)
}

func GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	return repository.GetEmployeeAttendances(employeeID, startDate, endDate)
}

// CorrectionRequest defines the payload for a manual attendance correction.
type CorrectionRequest struct {
	EmployeeID     int       `json:"employee_id" binding:"required"`
	CorrectionTime time.Time `json:"correction_time" binding:"required"`
	CorrectionType string    `json:"correction_type" binding:"required,oneof=check_in check_out"`
	Notes          string    `json:"notes" binding:"required,min=10"`
}

// CorrectAttendance handles the business logic for manual attendance correction by an admin.
// CorrectAttendance handles the business logic for manual attendance correction by an admin.
func CorrectAttendance(adminID uint, req CorrectionRequest) (*models.AttendancesTable, error) {
	// 1. Find the employee
	employee, err := repository.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, fmt.Errorf("employee with ID %d not found", req.EmployeeID)
	}

	// 2. Handle based on correction type
	if req.CorrectionType == "check_out" {
		// Find the latest attendance record for that day that needs a check-out
		latestAttendance, err := repository.GetLatestAttendanceForDate(req.EmployeeID, req.CorrectionTime)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve attendance record for correction: %w", err)
		}

		if latestAttendance == nil || latestAttendance.CheckOutTime != nil {
			return nil, fmt.Errorf("no pending check-in found for this employee on the selected date to apply a check-out")
		}

		// Update the existing record
		now := req.CorrectionTime
		latestAttendance.CheckOutTime = &now
		latestAttendance.Status = "present (corrected)"
		latestAttendance.IsCorrection = true
		latestAttendance.Notes = req.Notes
		latestAttendance.CorrectedByAdminID = &adminID

		if err := repository.UpdateAttendance(latestAttendance); err != nil {
			return nil, fmt.Errorf("failed to save corrected attendance: %w", err)
		}
		return latestAttendance, nil

	} else if req.CorrectionType == "check_in" {
		// Create a new attendance record because admin is manually adding a full day's record (or just a check-in)
		newAttendance := &models.AttendancesTable{
			EmployeeID:         req.EmployeeID,
			CheckInTime:        req.CorrectionTime,
			Status:             "present (corrected)",
			IsCorrection:       true,
			Notes:              req.Notes,
			CorrectedByAdminID: &adminID,
		}

		if err := repository.CreateAttendance(newAttendance); err != nil {
			return nil, fmt.Errorf("failed to create new corrected attendance: %w", err)
		}
		return newAttendance, nil
	}

	return nil, fmt.Errorf("invalid correction type specified")
}

// MarkDailyAbsentees checks for employees who haven't checked in and aren't on leave, and marks them as absent.
func MarkDailyAbsentees() error {
	log.Println("Starting daily absentee marking process...")

	companies, err := repository.GetAllActiveCompanies()
	if err != nil {
		return fmt.Errorf("failed to get active companies: %w", err)
	}

	for _, company := range companies {
		log.Printf("Processing company: %s (ID: %d)", company.Name, company.ID)

		companyLocation, err := time.LoadLocation(company.Timezone)
		if err != nil {
			log.Printf("Error loading company timezone %s for company %d: %v", company.Timezone, company.ID, err)
			continue // Skip this company if timezone is invalid
		}

		employees, err := repository.GetActiveEmployeesByCompanyID(company.ID)
		if err != nil {
			log.Printf("Failed to get active employees for company %d: %v", company.ID, err)
			continue
		}

		shifts, err := repository.GetShiftsByCompanyID(company.ID)
		if err != nil {
			log.Printf("Failed to get shifts for company %d: %v", company.ID, err)
			continue
		}

		// Create a map for quick shift lookup by ID
		shiftMap := make(map[uint]models.ShiftsTable)
		for _, shift := range shifts {
			shiftMap[uint(shift.ID)] = shift
		}

		for _, employee := range employees {
			// Skip if employee has no shift assigned
			if employee.ShiftID == nil {
				log.Printf("Employee %s (ID: %d) has no shift assigned. Skipping.", employee.Name, employee.ID)
				continue
			}

			shift, ok := shiftMap[uint(*employee.ShiftID)]
			if !ok {
				log.Printf("Shift with ID %d not found for employee %s (ID: %d). Skipping.", *employee.ShiftID, employee.Name, employee.ID)
				continue
			}

			// Calculate the grace period end time for the shift
			shiftEnd, err := helper.ParseTime(time.Now().In(companyLocation), shift.EndTime, companyLocation)
			if err != nil {
				log.Printf("Error parsing shift end time %s for employee %s (ID: %d): %v", shift.EndTime, employee.Name, employee.ID, err)
				continue
			}

			// If shift crosses midnight, adjust shiftEnd to be on the next day
			shiftStart, err := helper.ParseTime(time.Now().In(companyLocation), shift.StartTime, companyLocation)
			if err != nil {
				log.Printf("Error parsing shift start time %s for employee %s (ID: %d): %v", shift.StartTime, employee.Name, employee.ID, err)
				continue
			}
			if shiftEnd.Before(shiftStart) {
				shiftEnd = shiftEnd.Add(24 * time.Hour)
			}

			// Define the grace period after shift ends (e.g., 5 hours)
			gracePeriodAfterShift := 5 * time.Hour // This can be configurable
			processingCutoffTime := shiftEnd.Add(gracePeriodAfterShift)

			// Only process if the current time is past the processing cutoff time
			nowInCompanyLocation := time.Now().In(companyLocation)
			if nowInCompanyLocation.Before(processingCutoffTime) {
				log.Printf("Current time %s is before processing cutoff %s for employee %s (ID: %d). Skipping.", nowInCompanyLocation.Format("15:04"), processingCutoffTime.Format("15:04"), employee.Name, employee.ID)
				continue
			}

			// Check if already has an attendance record for today
			hasAttendance, err := repository.HasAttendanceForDate(employee.ID, time.Now().In(companyLocation))
			if err != nil {
				log.Printf("Error checking attendance for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
				continue
			}

			if hasAttendance {
				log.Printf("Employee %s (ID: %d) already has attendance for today. Skipping.", employee.Name, employee.ID)
				continue
			}

			// Check if employee is on approved leave for today
			approvedLeave, err := repository.IsEmployeeOnApprovedLeave(employee.ID, time.Now().In(companyLocation))
			if err != nil {
				log.Printf("Error checking leave status for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
				continue
			}

			if approvedLeave != nil {
				var status string
				var notes string
				if approvedLeave.Type == "sakit" {
					status = "on_sick"
					notes = "Automatically marked as on sick leave due to approved sick request."
				} else {
					status = "on_leave"
					notes = "Automatically marked as on leave due to approved leave request."
				}
				log.Printf("Employee %s (ID: %d) is on approved %s. Creating '%s' record.", employee.Name, employee.ID, approvedLeave.Type, status)
				absenceTime := time.Now().In(companyLocation) // Record the time of marking
				newAttendance := &models.AttendancesTable{
					EmployeeID:  employee.ID,
					CheckInTime: absenceTime,
					Status:      status,
					IsCorrection: true, // Mark as correction as it's not a physical check-in
					Notes:       notes,
				}
				if err := repository.CreateAttendance(newAttendance); err != nil {
					log.Printf("Failed to create %s record for employee %s (ID: %d): %v", status, employee.Name, employee.ID, err)
				}
				continue // Move to next employee after marking as on_leave/on_sick
			}

			// If no attendance and not on leave, mark as absent
			log.Printf("Marking employee %s (ID: %d) as absent for today.", employee.Name, employee.ID)
			absenceTime := time.Now().In(companyLocation) // Record the time of marking
			newAttendance := &models.AttendancesTable{
				EmployeeID:  employee.ID,
				CheckInTime: absenceTime,
				Status:      "absent",
				IsCorrection: true,
				Notes:       "Automatically marked as absent due to no check-in and no approved leave.",
			}
			if err := repository.CreateAttendance(newAttendance); err != nil {
				log.Printf("Failed to create absent record for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
			}
		}
	}

	log.Println("Daily absentee marking process finished.")
	return nil
}


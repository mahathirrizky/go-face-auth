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

	now := time.Now().In(companyLocation) // Get current time in company's timezone
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
		// Check if current time is within regular shift
		isWithinShift, err := helper.IsTimeWithinShift(now, shift.StartTime, shift.EndTime, shift.GracePeriodMinutes, companyLocation)
		if err != nil {
			log.Printf("Error checking time within shift: %v", err)
			return "", nil, time.Time{}, fmt.Errorf("failed to validate shift time")
		}

		if !isWithinShift {
			return "", nil, time.Time{}, fmt.Errorf("cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead")
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

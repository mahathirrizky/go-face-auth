package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

type AttendanceService interface {
	HandleAttendance(req AttendanceRequest) (string, *models.EmployeesTable, time.Time, error)
	HandleOvertimeCheckIn(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, error)
	HandleOvertimeCheckOut(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, int, time.Time, error)
	GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error)
	ExportEmployeeAttendanceToExcel(employeeID int, startDate, endDate *time.Time) (*excelize.File, string, error)
	ExportAllAttendancesToExcel(companyID int, startDate, endDate *time.Time) (*excelize.File, string, error)
	ExportUnaccountedToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error)
	ExportOvertimeToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error)
	GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error)
	GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error)
	GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error)
	CorrectAttendance(adminID uint, req CorrectionRequest) (*models.AttendancesTable, error)
	MarkDailyAbsentees() error
}

type attendanceService struct {
	employeeRepo    repository.EmployeeRepository
	companyRepo     repository.CompanyRepository
	attendanceRepo  repository.AttendanceRepository
	faceImageRepo   repository.FaceImageRepository
	locationRepo    repository.AttendanceLocationRepository
	leaveRequestRepo repository.LeaveRequestRepository
	shiftRepo       repository.ShiftRepository
	divisionRepo    repository.DivisionRepository // Added divisionRepo
	pythonClient    PythonServerClientInterface
}

func NewAttendanceService(employeeRepo repository.EmployeeRepository, companyRepo repository.CompanyRepository, attendanceRepo repository.AttendanceRepository, faceImageRepo repository.FaceImageRepository, locationRepo repository.AttendanceLocationRepository, leaveRequestRepo repository.LeaveRequestRepository, shiftRepo repository.ShiftRepository, divisionRepo repository.DivisionRepository, pythonClient PythonServerClientInterface) AttendanceService {
	return &attendanceService{
		employeeRepo:    employeeRepo,
		companyRepo:     companyRepo,
		attendanceRepo:  attendanceRepo,
		faceImageRepo:   faceImageRepo,
		locationRepo:    locationRepo,
		leaveRequestRepo: leaveRequestRepo,
		shiftRepo:       shiftRepo,
		divisionRepo:    divisionRepo, // Added divisionRepo
		pythonClient:    pythonClient,
	}
}

// AttendanceRequest represents the request body for attendance.
type AttendanceRequest struct {
	EmployeeID int     `json:"employee_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	ImageData  string  `json:"image_data" binding:"required"`
}




	func (s *attendanceService) HandleAttendance(req AttendanceRequest) (string, *models.EmployeesTable, time.Time, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return "", nil, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := s.companyRepo.GetCompanyByID(employee.CompanyID)
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
	approvedLeave, err := s.leaveRequestRepo.IsEmployeeOnApprovedLeave(employee.ID, now)
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
	faceImages, err := s.faceImageRepo.GetFaceImagesByEmployeeID(req.EmployeeID)
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

	pythonResponse, err := s.pythonClient.SendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return "", nil, time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return "", nil, time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	// Determine effective shift and locations
	var effectiveShift models.ShiftsTable
	var effectiveLocations []models.AttendanceLocation

	if employee.DivisionID != nil {
		division, err := s.divisionRepo.GetDivisionByID(uint(*employee.DivisionID))
		if err == nil && division != nil {
			if len(division.Shifts) > 0 {
				effectiveShift = division.Shifts[0] // Assuming one shift per division for simplicity, or pick default
			} else if employee.ShiftID != nil {
				effectiveShift = employee.Shift // Fallback to employee's assigned shift
			}
			if len(division.Locations) > 0 {
				effectiveLocations = division.Locations
			} else {
				effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
				if err != nil {
					return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
				}
			}
		} else {
			// Division not found or error, fallback to employee's assigned
			if employee.ShiftID != nil {
				effectiveShift = employee.Shift
			}
			effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
			if err != nil {
				return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
			}
		}
	} else {
		// No division assigned, use employee's assigned shift and company locations
		if employee.ShiftID != nil {
			effectiveShift = employee.Shift
		}
		effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
		if err != nil {
			return "", nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
		}
	}

	// Validate effective shift and locations
	if effectiveShift.ID == 0 {
		return "", nil, time.Time{}, fmt.Errorf("employee does not have an assigned shift or division shift")
	}
	if len(effectiveLocations) == 0 {
		return "", nil, time.Time{}, fmt.Errorf("no valid attendance locations configured for employee or division")
	}

	// Validate employee's current location against effective attendance locations
	isWithinValidLocation := false
	for _, loc := range effectiveLocations {
		distance := helper.HaversineDistance(req.Latitude, req.Longitude, loc.Latitude, loc.Longitude)
		if distance <= float64(loc.Radius) {
			isWithinValidLocation = true
			break
		}
	}

	if !isWithinValidLocation {
		return "", nil, time.Time{}, fmt.Errorf("you are not within a valid attendance location")
	}

	var message string

	todaysAttendance, err := s.attendanceRepo.GetLatestAttendanceForDate(req.EmployeeID, now)
	if err != nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to retrieve attendance record for today")
	}

	if todaysAttendance == nil {
		// CASE 1: NO ATTENDANCE RECORD FOR TODAY - THIS IS A CHECK-IN
		// Calculate earliest allowed check-in time (1.5 hours before shift start)
		shiftStartToday, err := helper.ParseTime(now, effectiveShift.StartTime, companyLocation)
		if err != nil {
			log.Printf("Error parsing shift start time for early check-in: %v", err)
			return "", nil, time.Time{}, fmt.Errorf("failed to validate shift time")
		}
		earliesCheckInTime := shiftStartToday.Add(-90 * time.Minute) // 90 minutes = 1.5 hours

		// Prevent check-in if too early
		if now.Before(earliesCheckInTime) {
			return "", nil, time.Time{}, fmt.Errorf("anda tidak dapat absen lebih dari 1.5 jam sebelum jam shift Anda")
		}

		// Check if current time is within regular shift (considering grace period for late check-in)
		isWithinShift, err := helper.IsTimeWithinShift(now, effectiveShift.StartTime, effectiveShift.EndTime, effectiveShift.GracePeriodMinutes, companyLocation)
		if err != nil {
			log.Printf("Error checking time within shift: %v", err)
			return "", nil, time.Time{}, fmt.Errorf("failed to validate shift time")
		}

		if !isWithinShift {
			return "", nil, time.Time{}, fmt.Errorf("cannot check-in for regular attendance outside of shift hours. Use overtime check-in instead")
		}

		if now.After(shiftStartToday.Add(time.Duration(effectiveShift.GracePeriodMinutes) * time.Minute)) {
			status = "late"
		} else {
			status = "on_time"
		}

		newAttendance := &models.AttendancesTable{
			EmployeeID:  req.EmployeeID,
			CheckInTime: now,
			Status:      status,
		}
		err = s.attendanceRepo.CreateAttendance(newAttendance)
		message = "Check-in successful!"

	} else if todaysAttendance.CheckOutTime == nil {
		// CASE 2: ATTENDANCE RECORD EXISTS BUT NO CHECK-OUT - THIS IS A CHECK-OUT
		todaysAttendance.CheckOutTime = &now
		todaysAttendance.Status = "present"
		err = s.attendanceRepo.UpdateAttendance(todaysAttendance)
		message = "Check-out successful!"

	} else {
		// CASE 3: ATTENDANCE RECORD EXISTS AND IS ALREADY CHECKED OUT
		return "", nil, time.Time{}, fmt.Errorf("anda sudah melakukan check-in dan check-out untuk hari ini")
	}

	if err != nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to record attendance: %w", err)
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
func (s *attendanceService) HandleOvertimeCheckIn(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, error) {
	// --- Face Recognition Logic ---
	faceImages, err := s.faceImageRepo.GetFaceImagesByEmployeeID(req.EmployeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", req.EmployeeID, err)
		return nil, time.Time{}, fmt.Errorf("could not retrieve employee face image")
	}
	if len(faceImages) == 0 {
		return nil, time.Time{}, fmt.Errorf("no registered face images for this employee")
	}
	dbImagePath := faceImages[0].ImagePath

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: req.ImageData,
		DBImagePath:     dbImagePath,
	}

	pythonResponse, err := s.pythonClient.SendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return nil, time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return nil, time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := s.companyRepo.GetCompanyByID(employee.CompanyID)
	if err != nil || company == nil {
		return nil, time.Time{}, fmt.Errorf("failed to retrieve company information")
	}

	companyLocation, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return nil, time.Time{}, fmt.Errorf("invalid company timezone configuration")
	}

	// Determine effective shift and locations
	var effectiveShift models.ShiftsTable
	var effectiveLocations []models.AttendanceLocation

	if employee.DivisionID != nil {
		division, err := s.divisionRepo.GetDivisionByID(uint(*employee.DivisionID))
		if err == nil && division != nil {
			if len(division.Shifts) > 0 {
				effectiveShift = division.Shifts[0] // Assuming one shift per division for simplicity, or pick default
			} else if employee.ShiftID != nil {
				effectiveShift = employee.Shift // Fallback to employee's assigned shift
			}
			if len(division.Locations) > 0 {
				effectiveLocations = division.Locations
			} else {
				effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
				if err != nil {
					return nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
				}
			}
		} else {
			// Division not found or error, fallback to employee's assigned
			if employee.ShiftID != nil {
				effectiveShift = employee.Shift
			}
			effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
			if err != nil {
				return nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
			}
		}
	} else {
		// No division assigned, use employee's assigned shift and company locations
		if employee.ShiftID != nil {
			effectiveShift = employee.Shift
		}
		effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to retrieve company attendance locations")
		}
	}

	// Validate effective shift and locations
	if effectiveShift.ID == 0 {
		return nil, time.Time{}, fmt.Errorf("employee does not have an assigned shift or division shift")
	}
	if len(effectiveLocations) == 0 {
		return nil, time.Time{}, fmt.Errorf("no valid attendance locations configured for employee or division")
	}

	// Validate employee's current location against effective attendance locations
	isWithinValidLocation := false
	for _, loc := range effectiveLocations {
		distance := helper.HaversineDistance(req.Latitude, req.Longitude, loc.Latitude, loc.Longitude)
		if distance <= float64(loc.Radius) {
			isWithinValidLocation = true
			break
		}
	}

	if !isWithinValidLocation {
		return nil, time.Time{}, fmt.Errorf("you are not within a valid attendance location")
	}

	// Get employee's shift
	// if employee.ShiftID == nil {
	// 	return nil, time.Time{}, fmt.Errorf("employee does not have a shift assigned")
	// }
	// shift := employee.Shift

	now := time.Now().In(companyLocation) // Get current time in company's timezone

	// Validate: Cannot check-in for overtime if within regular shift hours
	isWithinShift, err := helper.IsTimeWithinShift(now, effectiveShift.StartTime, effectiveShift.EndTime, effectiveShift.GracePeriodMinutes, companyLocation)
	if err != nil {
		log.Printf("Error checking time within shift for overtime check-in: %v", err)
		return nil, time.Time{}, fmt.Errorf("failed to validate shift time")
	}
	if isWithinShift {
		return nil, time.Time{}, fmt.Errorf("cannot check-in for overtime during regular shift hours")
	}

	// Check if employee has an open regular check-in
	latestRegularAttendance, err := s.attendanceRepo.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to retrieve latest regular attendance record")
	}
	if latestRegularAttendance != nil && latestRegularAttendance.CheckOutTime == nil && latestRegularAttendance.Status != "overtime_in" && latestRegularAttendance.Status != "overtime_out" {
		return nil, time.Time{}, fmt.Errorf("anda harus check-out dari shift reguler sebelum check-in lembur")
	}

	// Check if employee is already checked in for overtime
	latestOvertimeAttendance, err := s.attendanceRepo.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to retrieve latest overtime record")
	}
	if latestOvertimeAttendance != nil && latestOvertimeAttendance.CheckOutTime == nil && latestOvertimeAttendance.Status == "overtime_in" {
		return nil, time.Time{}, fmt.Errorf("employee is already checked in for overtime")
	}

	// Create new overtime check-in record
	newOvertimeAttendance := &models.AttendancesTable{
		EmployeeID:  req.EmployeeID,
		CheckInTime: now,
		Status:      "overtime_in", // Specific status for overtime check-in
	}
	err = s.attendanceRepo.CreateAttendance(newOvertimeAttendance)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to record overtime check-in")
	}

	return employee, now, nil
}

func (s *attendanceService) HandleOvertimeCheckOut(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, int, time.Time, error) {
	// --- Face Recognition Logic ---
	faceImages, err := s.faceImageRepo.GetFaceImagesByEmployeeID(req.EmployeeID)
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

	pythonResponse, err := s.pythonClient.SendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("face recognition service is unavailable")
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("face not recognized")
	}
	// --- End of Face Recognition Logic ---

	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("employee not found")
	}

	// Get employee's company and its timezone
	company, err := s.companyRepo.GetCompanyByID(employee.CompanyID)
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
	latestOvertimeAttendance, err := s.attendanceRepo.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
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

	err = s.attendanceRepo.UpdateAttendance(latestOvertimeAttendance)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("failed to record overtime check-out")
	}

	// Return employee, now, overtimeMinutes, and original CheckInTime
	return employee, now, overtimeMinutes, latestOvertimeAttendance.CheckInTime, nil
}

func (s *attendanceService) GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error) {
	return s.attendanceRepo.GetAttendancesPaginated(companyID, startDate, endDate, search, page, pageSize)
}

func (s *attendanceService) ExportEmployeeAttendanceToExcel(employeeID int, startDate, endDate *time.Time) (*excelize.File, string, error) {
	attendances, err := s.attendanceRepo.GetEmployeeAttendances(employeeID, startDate, endDate)
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

func (s *attendanceService) ExportAllAttendancesToExcel(companyID int, startDate, endDate *time.Time) (*excelize.File, string, error) {
	attendances, err := s.attendanceRepo.GetCompanyAttendancesFiltered(companyID, startDate, endDate, "all")
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

func (s *attendanceService) ExportUnaccountedToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error) {
	unaccountedEmployees, err := s.attendanceRepo.GetUnaccountedEmployeesFiltered(companyID, startDate, endDate, search)
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

func (s *attendanceService) ExportOvertimeToExcel(companyID int, startDate, endDate *time.Time, search string) (*excelize.File, string, error) {
	overtimeAttendances, err := s.attendanceRepo.GetOvertimeAttendancesFiltered(companyID, startDate, endDate, search)
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

func (s *attendanceService) GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.AttendancesTable, int64, error) {
	return s.attendanceRepo.GetOvertimeAttendancesPaginated(companyID, startDate, endDate, search, page, pageSize)
}
func (s *attendanceService) GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
	return s.attendanceRepo.GetUnaccountedEmployeesPaginated(companyID, startDate, endDate, search, page, pageSize)
}

func (s *attendanceService) GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	return s.attendanceRepo.GetEmployeeAttendances(employeeID, startDate, endDate)
}

// CorrectionRequest defines the payload for a manual attendance correction.
type CorrectionRequest struct {
	EmployeeID     int       `json:"employee_id" binding:"required"`
	CorrectionTime time.Time `json:"correction_time" binding:"required"`
	CorrectionType string    `json:"correction_type" binding:"required,oneof=check_in check_out"`
	Notes          string    `json:"notes" binding:"required,min=10"`
}

// CorrectAttendance handles the business logic for manual attendance correction by an admin.
func (s *attendanceService) CorrectAttendance(adminID uint, req CorrectionRequest) (*models.AttendancesTable, error) {
	// 1. Find the employee
	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, fmt.Errorf("employee with ID %d not found", req.EmployeeID)
	}

	// 2. Handle based on correction type
	if req.CorrectionType == "check_out" {
		// Find the latest attendance record for that day that needs a check-out
		latestAttendance, err := s.attendanceRepo.GetLatestAttendanceForDate(req.EmployeeID, req.CorrectionTime)
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

		if err := s.attendanceRepo.UpdateAttendance(latestAttendance); err != nil {
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

		if err := s.attendanceRepo.CreateAttendance(newAttendance); err != nil {
			return nil, fmt.Errorf("failed to create new corrected attendance: %w", err)
		}
		return newAttendance, nil
	}

	return nil, fmt.Errorf("invalid correction type specified")
}

// MarkDailyAbsentees checks for employees who haven't checked in and aren't on leave, and marks them as absent.
// It also cleans up incomplete attendance records from the previous day.
func (s *attendanceService) MarkDailyAbsentees() error {
	log.Println("Starting daily absentee and cleanup process...")

	companies, err := s.companyRepo.GetAllActiveCompanies()
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

		nowInCompanyLocation := time.Now().In(companyLocation)
		yesterday := nowInCompanyLocation.AddDate(0, 0, -1)

		// --- Cleanup: Mark incomplete attendances from yesterday ---
		log.Printf("Checking for incomplete attendances from %s for company %s", yesterday.Format("2006-01-02"), company.Name)
		incompleteAttendances, err := s.attendanceRepo.FindIncompleteAttendancesByCompany(company.ID, yesterday)
		if err != nil {
			log.Printf("Error finding incomplete attendances for company %d: %v", company.ID, err)
		} else if len(incompleteAttendances) > 0 {
			log.Printf("Found %d incomplete attendance records to clean up.", len(incompleteAttendances))
			for _, att := range incompleteAttendances {
				attToUpdate := att // Make a new variable to avoid loop variable issues
				attToUpdate.Status = "incomplete"
				attToUpdate.Notes = "Automatically marked due to forgotten check-out."
				attToUpdate.IsCorrection = true
				if err := s.attendanceRepo.UpdateAttendance(&attToUpdate); err != nil {
					log.Printf("Failed to update incomplete attendance record %d: %v", attToUpdate.ID, err)
				} else {
					log.Printf("Marked attendance record %d as incomplete.", attToUpdate.ID)
				}
			}
		}
		// --- End of Cleanup ---

		employees, err := s.employeeRepo.GetActiveEmployeesByCompanyID(company.ID)
		if err != nil {
			log.Printf("Failed to get active employees for company %d: %v", company.ID, err)
			continue
		}

		shifts, err := s.shiftRepo.GetShiftsByCompanyID(company.ID)
		if err != nil {
			log.Printf("Failed to get shifts for company %d: %v", company.ID, err)
			continue
		}

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

			shiftEnd, err := helper.ParseTime(nowInCompanyLocation, shift.EndTime, companyLocation)
			if err != nil {
				log.Printf("Error parsing shift end time %s for employee %s (ID: %d): %v", shift.EndTime, employee.Name, employee.ID, err)
				continue
			}

			shiftStart, err := helper.ParseTime(nowInCompanyLocation, shift.StartTime, companyLocation)
			if err != nil {
				log.Printf("Error parsing shift start time %s for employee %s (ID: %d): %v", shift.StartTime, employee.Name, employee.ID, err)
				continue
			}
			if shiftEnd.Before(shiftStart) {
				shiftEnd = shiftEnd.Add(24 * time.Hour)
			}

			gracePeriodAfterShift := 5 * time.Hour
			processingCutoffTime := shiftEnd.Add(gracePeriodAfterShift)

			if nowInCompanyLocation.Before(processingCutoffTime) {
				log.Printf("Current time %s is before processing cutoff %s for employee %s (ID: %d). Skipping.", nowInCompanyLocation.Format("15:04"), processingCutoffTime.Format("15:04"), employee.Name, employee.ID)
				continue
			}

			hasAttendance, err := s.attendanceRepo.HasAttendanceForDate(employee.ID, nowInCompanyLocation)
			if err != nil {
				log.Printf("Error checking attendance for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
				continue
			}

			if hasAttendance {
				log.Printf("Employee %s (ID: %d) already has attendance for today. Skipping.", employee.Name, employee.ID)
				continue
			}

			approvedLeave, err := s.leaveRequestRepo.IsEmployeeOnApprovedLeave(employee.ID, nowInCompanyLocation)
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
				absenceTime := nowInCompanyLocation
				newAttendance := &models.AttendancesTable{
					EmployeeID:   employee.ID,
					CheckInTime:  absenceTime,
					Status:       status,
					IsCorrection: true,
					Notes:        notes,
				}
				if err := s.attendanceRepo.CreateAttendance(newAttendance); err != nil {
					log.Printf("Failed to create %s record for employee %s (ID: %d): %v", status, employee.Name, employee.ID, err)
				}
				continue
			}

			log.Printf("Marking employee %s (ID: %d) as absent for today.", employee.Name, employee.ID)
			absenceTime := nowInCompanyLocation
			newAttendance := &models.AttendancesTable{
				EmployeeID:   employee.ID,
				CheckInTime:  absenceTime,
				Status:       "absent",
				IsCorrection: true,
				Notes:        "Automatically marked as absent due to no check-in and no approved leave.",
			}
			if err := s.attendanceRepo.CreateAttendance(newAttendance); err != nil {
				log.Printf("Failed to create absent record for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
			}
		}
	}

	log.Println("Daily absentee and cleanup process finished.")
	return nil
}
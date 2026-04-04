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

// Constants for attendance business rules
const (
	// EarlyCheckInWindow is how early before shift start an employee can check in.
	EarlyCheckInWindow = 90 * time.Minute // 1.5 hours

	// GracePeriodAfterShift is the buffer after shift end before marking absent.
	GracePeriodAfterShift = 5 * time.Hour
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
	employeeRepo     repository.EmployeeRepository
	companyRepo      repository.CompanyRepository
	attendanceRepo   repository.AttendanceRepository
	faceImageRepo    repository.FaceImageRepository
	locationRepo     repository.AttendanceLocationRepository
	leaveRequestRepo repository.LeaveRequestRepository
	shiftRepo        repository.ShiftRepository
	divisionRepo     repository.DivisionRepository
	pythonClient     PythonServerClientInterface
}

func NewAttendanceService(employeeRepo repository.EmployeeRepository, companyRepo repository.CompanyRepository, attendanceRepo repository.AttendanceRepository, faceImageRepo repository.FaceImageRepository, locationRepo repository.AttendanceLocationRepository, leaveRequestRepo repository.LeaveRequestRepository, shiftRepo repository.ShiftRepository, divisionRepo repository.DivisionRepository, pythonClient PythonServerClientInterface) AttendanceService {
	return &attendanceService{
		employeeRepo:     employeeRepo,
		companyRepo:      companyRepo,
		attendanceRepo:   attendanceRepo,
		faceImageRepo:    faceImageRepo,
		locationRepo:     locationRepo,
		leaveRequestRepo: leaveRequestRepo,
		shiftRepo:        shiftRepo,
		divisionRepo:     divisionRepo,
		pythonClient:     pythonClient,
	}
}

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

// --- Private helper methods to eliminate code duplication ---

// verifyFaceRecognition performs face recognition against the employee's registered face images.
func (s *attendanceService) verifyFaceRecognition(employeeID int, imageData string) error {
	faceImages, err := s.faceImageRepo.GetFaceImagesByEmployeeID(employeeID)
	if err != nil {
		log.Printf("Error getting face image from DB for employee %d: %v", employeeID, err)
		return ErrFaceImageRetrieval
	}
	if len(faceImages) == 0 {
		return ErrNoRegisteredFaceImages
	}

	pythonPayload := PythonRecognitionRequest{
		ClientImageData: imageData,
		DBImagePath:     faceImages[0].ImagePath,
	}

	pythonResponse, err := s.pythonClient.SendToPythonServer(pythonPayload)
	if err != nil {
		log.Printf("Error communicating with Python server: %v", err)
		return ErrFaceRecognitionUnavailable
	}

	status, ok := pythonResponse["status"].(string)
	if !ok || status != "recognized" {
		return ErrFaceNotRecognized
	}

	return nil
}

// getCompanyTimezone loads the timezone for a given company.
func (s *attendanceService) getCompanyTimezone(companyID int) (*time.Location, *models.CompaniesTable, error) {
	company, err := s.companyRepo.GetCompanyByID(companyID)
	if err != nil || company == nil {
		return nil, nil, ErrCompanyNotFound
	}

	loc, err := time.LoadLocation(company.Timezone)
	if err != nil {
		log.Printf("Error loading company timezone %s: %v", company.Timezone, err)
		return nil, nil, ErrInvalidTimezone
	}

	return loc, company, nil
}

// resolveEffectiveShiftAndLocations determines the effective shift and attendance locations
// for an employee based on their division assignment (if any) or direct assignment.
func (s *attendanceService) resolveEffectiveShiftAndLocations(employee *models.EmployeesTable) (models.ShiftsTable, []models.AttendanceLocation, error) {
	var effectiveShift models.ShiftsTable
	var effectiveLocations []models.AttendanceLocation

	if employee.DivisionID != nil {
		division, err := s.divisionRepo.GetDivisionByID(uint(*employee.DivisionID))
		if err == nil && division != nil {
			if len(division.Shifts) > 0 {
				effectiveShift = division.Shifts[0]
			} else if employee.ShiftID != nil {
				effectiveShift = employee.Shift
			}
			if len(division.Locations) > 0 {
				effectiveLocations = division.Locations
			} else {
				effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
				if err != nil {
					return effectiveShift, nil, ErrLocationRetrieval
				}
			}
		} else {
			// Division not found or error, fallback to employee's assigned
			if employee.ShiftID != nil {
				effectiveShift = employee.Shift
			}
			effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
			if err != nil {
				return effectiveShift, nil, ErrLocationRetrieval
			}
		}
	} else {
		// No division assigned, use employee's assigned shift and company locations
		if employee.ShiftID != nil {
			effectiveShift = employee.Shift
		}
		var err error
		effectiveLocations, err = s.locationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
		if err != nil {
			return effectiveShift, nil, ErrLocationRetrieval
		}
	}

	// Validate
	if effectiveShift.ID == 0 {
		return effectiveShift, nil, ErrNoShiftAssigned
	}
	if len(effectiveLocations) == 0 {
		return effectiveShift, nil, ErrNoLocationsConfigured
	}

	return effectiveShift, effectiveLocations, nil
}

// validateLocation checks if the given coordinates are within any of the attendance locations.
func (s *attendanceService) validateLocation(latitude, longitude float64, locations []models.AttendanceLocation) error {
	for _, loc := range locations {
		distance := helper.HaversineDistance(latitude, longitude, loc.Latitude, loc.Longitude)
		if distance <= float64(loc.Radius) {
			return nil
		}
	}
	return ErrOutsideAttendanceLocation
}

// --- Main attendance handlers ---

func (s *attendanceService) HandleAttendance(req AttendanceRequest) (string, *models.EmployeesTable, time.Time, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return "", nil, time.Time{}, ErrEmployeeNotFound
	}

	companyLocation, _, err := s.getCompanyTimezone(employee.CompanyID)
	if err != nil {
		return "", nil, time.Time{}, err
	}

	now := time.Now().In(companyLocation)

	// Check leave status
	approvedLeave, err := s.leaveRequestRepo.IsEmployeeOnApprovedLeave(employee.ID, now)
	if err != nil {
		log.Printf("Error checking leave status for employee %s (ID: %d): %v", employee.Name, employee.ID, err)
		return "", nil, time.Time{}, ErrLeaveCheckFailed
	}
	if approvedLeave != nil {
		leaveType := "cuti"
		if approvedLeave.Type == "sakit" {
			leaveType = "sakit"
		}
		return "", nil, time.Time{}, fmt.Errorf("anda sedang dalam pengajuan %s yang disetujui untuk hari ini", leaveType)
	}

	// Face recognition
	if err := s.verifyFaceRecognition(req.EmployeeID, req.ImageData); err != nil {
		return "", nil, time.Time{}, err
	}

	// Resolve shift and locations
	effectiveShift, effectiveLocations, err := s.resolveEffectiveShiftAndLocations(employee)
	if err != nil {
		return "", nil, time.Time{}, err
	}

	// Validate location
	if err := s.validateLocation(req.Latitude, req.Longitude, effectiveLocations); err != nil {
		return "", nil, time.Time{}, err
	}

	var message string
	var status string

	todaysAttendance, err := s.attendanceRepo.GetLatestAttendanceForDate(req.EmployeeID, now)
	if err != nil {
		return "", nil, time.Time{}, ErrAttendanceRetrieval
	}

	if todaysAttendance == nil {
		// CASE 1: CHECK-IN
		shiftStartToday, err := helper.ParseTime(now, effectiveShift.StartTime, companyLocation)
		if err != nil {
			log.Printf("Error parsing shift start time for early check-in: %v", err)
			return "", nil, time.Time{}, ErrShiftValidationFailed
		}
		earliestCheckInTime := shiftStartToday.Add(-EarlyCheckInWindow)

		if now.Before(earliestCheckInTime) {
			return "", nil, time.Time{}, ErrTooEarlyForCheckIn
		}

		isWithinShift, err := helper.IsTimeWithinShift(now, effectiveShift.StartTime, effectiveShift.EndTime, effectiveShift.GracePeriodMinutes, companyLocation)
		if err != nil {
			log.Printf("Error checking time within shift: %v", err)
			return "", nil, time.Time{}, ErrShiftValidationFailed
		}
		if !isWithinShift {
			return "", nil, time.Time{}, ErrOutsideShiftHours
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
		// CASE 2: CHECK-OUT
		todaysAttendance.CheckOutTime = &now
		todaysAttendance.Status = "present"
		err = s.attendanceRepo.UpdateAttendance(todaysAttendance)
		message = "Check-out successful!"

	} else {
		// CASE 3: ALREADY DONE
		return "", nil, time.Time{}, ErrAlreadyCheckedOut
	}

	if err != nil {
		return "", nil, time.Time{}, fmt.Errorf("failed to record attendance: %w", err)
	}

	return message, employee, now, nil
}


func (s *attendanceService) HandleOvertimeCheckIn(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, time.Time{}, ErrEmployeeNotFound
	}

	companyLocation, _, err := s.getCompanyTimezone(employee.CompanyID)
	if err != nil {
		return nil, time.Time{}, err
	}

	now := time.Now().In(companyLocation)

	if err := s.verifyFaceRecognition(req.EmployeeID, req.ImageData); err != nil {
		return nil, time.Time{}, err
	}

	effectiveShift, effectiveLocations, err := s.resolveEffectiveShiftAndLocations(employee)
	if err != nil {
		return nil, time.Time{}, err
	}

	if err := s.validateLocation(req.Latitude, req.Longitude, effectiveLocations); err != nil {
		return nil, time.Time{}, err
	}

	// Validate: Cannot check-in for overtime if within regular shift hours
	isWithinShift, err := helper.IsTimeWithinShift(now, effectiveShift.StartTime, effectiveShift.EndTime, effectiveShift.GracePeriodMinutes, companyLocation)
	if err != nil {
		log.Printf("Error checking time within shift for overtime check-in: %v", err)
		return nil, time.Time{}, ErrShiftValidationFailed
	}
	if isWithinShift {
		return nil, time.Time{}, ErrOvertimeDuringShift
	}

	// Check if employee has an open regular check-in
	latestRegularAttendance, err := s.attendanceRepo.GetLatestAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, ErrAttendanceRetrieval
	}
	if latestRegularAttendance != nil && latestRegularAttendance.CheckOutTime == nil && latestRegularAttendance.Status != "overtime_in" && latestRegularAttendance.Status != "overtime_out" {
		return nil, time.Time{}, ErrMustCheckOutRegular
	}

	// Check if employee is already checked in for overtime
	latestOvertimeAttendance, err := s.attendanceRepo.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, ErrAttendanceRetrieval
	}
	if latestOvertimeAttendance != nil && latestOvertimeAttendance.CheckOutTime == nil && latestOvertimeAttendance.Status == "overtime_in" {
		return nil, time.Time{}, ErrAlreadyCheckedInOvertime
	}

	newOvertimeAttendance := &models.AttendancesTable{
		EmployeeID:  req.EmployeeID,
		CheckInTime: now,
		Status:      "overtime_in", // Specific status for overtime check-in
	}
	err = s.attendanceRepo.CreateAttendance(newOvertimeAttendance)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to record overtime check-in: %w", err)
	}

	return employee, now, nil
}

func (s *attendanceService) HandleOvertimeCheckOut(req OvertimeAttendanceRequest) (*models.EmployeesTable, time.Time, int, time.Time, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(req.EmployeeID)
	if err != nil || employee == nil {
		return nil, time.Time{}, 0, time.Time{}, ErrEmployeeNotFound
	}

	companyLocation, _, err := s.getCompanyTimezone(employee.CompanyID)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, err
	}

	now := time.Now().In(companyLocation)

	if err := s.verifyFaceRecognition(req.EmployeeID, req.ImageData); err != nil {
		return nil, time.Time{}, 0, time.Time{}, err
	}

	// Find the latest "overtime_in" record that is not checked out
	latestOvertimeAttendance, err := s.attendanceRepo.GetLatestOvertimeAttendanceByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, ErrAttendanceRetrieval
	}
	if latestOvertimeAttendance == nil || latestOvertimeAttendance.CheckOutTime != nil || latestOvertimeAttendance.Status != "overtime_in" {
		return nil, time.Time{}, 0, time.Time{}, ErrNotCheckedInForOvertime
	}

	overtimeDuration := now.Sub(latestOvertimeAttendance.CheckInTime)
	overtimeMinutes := int(overtimeDuration.Minutes())

	latestOvertimeAttendance.CheckOutTime = &now
	latestOvertimeAttendance.OvertimeMinutes = overtimeMinutes
	latestOvertimeAttendance.Status = "overtime_out" // Specific status for overtime check-out

	err = s.attendanceRepo.UpdateAttendance(latestOvertimeAttendance)
	if err != nil {
		return nil, time.Time{}, 0, time.Time{}, fmt.Errorf("failed to record overtime check-out: %w", err)
	}

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
	switch req.CorrectionType {
case "check_out":
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

	case "check_in":
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
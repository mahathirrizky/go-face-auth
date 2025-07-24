package services

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"io"
	"log"
	"mime/multipart"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// CreateEmployeeRequest defines the request structure for creating an employee.
type CreateEmployeeRequest struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Position         string `json:"position" binding:"required"`
	EmployeeIDNumber string `json:"employee_id_number" binding:"required"`
	ShiftID          *int   `json:"shift_id"`
}

// CreateEmployee handles the creation of a new employee, including subscription limit checks and initial password setup.
func CreateEmployee(ctx context.Context, companyID uint, req CreateEmployeeRequest) (*models.EmployeesTable, error) {
	// Retrieve company and its subscription package/custom offer
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").Preload("CustomOffer").First(&company, companyID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve company information: %w", err)
	}

	// Determine the effective MaxEmployees limit
	var maxEmployeesLimit int
	if company.CustomOfferID != nil && company.CustomOffer != nil {
		maxEmployeesLimit = company.CustomOffer.MaxEmployees
	} else if company.SubscriptionPackage.ID != 0 {
		maxEmployeesLimit = company.SubscriptionPackage.MaxEmployees
	} else {
		return nil, fmt.Errorf("company has no active subscription package or custom offer")
	}

	// Check current employee count
	var employeeCount int64
	if err := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", companyID).Count(&employeeCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count existing employees: %w", err)
	}

	// Check if adding a new employee would exceed the package limit
	if employeeCount >= int64(maxEmployeesLimit) {
		return nil, fmt.Errorf("employee limit reached for your current plan")
	}

	employee := &models.EmployeesTable{
		CompanyID:        int(companyID),
		Email:            req.Email,
		Name:             req.Name,
		Position:         req.Position,
		EmployeeIDNumber: req.EmployeeIDNumber,
		Role:             "employee", // Set default role to employee
	}

	// Determine the shift ID for the employee
	if req.ShiftID != nil {
		employee.ShiftID = req.ShiftID
	} else {
		// If no shift ID is provided, try to find the default shift for the company
		defaultShift, err := repository.GetDefaultShiftByCompanyID(int(companyID))
		if err != nil {
			log.Printf("No default shift found for company %d. Employee %s will be created without a shift: %v", companyID, req.Email, err)
		} else if defaultShift != nil {
			employee.ShiftID = &defaultShift.ID
		}
	}

	if err := repository.CreateEmployee(employee); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	// Generate password reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24) // Token valid for 24 hours

	passwordResetToken := &models.PasswordResetTokenTable{
		UserID:    employee.ID,
		TokenType: "employee_initial_password",
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := repository.CreatePasswordResetToken(passwordResetToken); err != nil {
		log.Printf("Error creating password reset token for employee %d: %v", employee.ID, err)
		return nil, fmt.Errorf("failed to generate initial password link: %w", err)
	}

	// Send email with password reset link
	resetLink := os.Getenv("FRONTEND_BASE_URL") + "/initial-password-setup?token=" + token
	if os.Getenv("FRONTEND_BASE_URL") == "" {
		resetLink = "http://localhost:5173/initial-password-setup?token=" + token // Fallback for development
	}

	// Run email sending in a goroutine to avoid blocking the main request
	go func(email, name, link string) {
		if err := helper.SendPasswordResetEmail(email, name, link); err != nil {
			log.Printf("Error sending initial password email to %s in background: %v", email, err)
		}
	}(employee.Email, employee.Name, resetLink)

	return employee, nil
}

func GetEmployeeByID(employeeID int, companyID uint) (*models.EmployeesTable, error) {
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}
	if employee == nil || employee.CompanyID != int(companyID) {
		return nil, nil // Let handler decide on 404
	}
	return employee, nil
}

func GetEmployeesByCompanyIDPaginated(companyID int, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
	return repository.GetEmployeesByCompanyIDPaginated(companyID, search, page, pageSize)
}

func SearchEmployees(companyID int, name string) ([]models.EmployeesTable, error) {
	return repository.SearchEmployees(companyID, name)
}

func UpdateEmployee(employeeID int, companyID uint, updates map[string]interface{}) error {
	// Verify employee belongs to this company
	existingEmployee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || existingEmployee == nil || existingEmployee.CompanyID != int(companyID) {
		return fmt.Errorf("employee not found or does not belong to your company")
	}

	// Prevent updating password via this generic update endpoint
	delete(updates, "password")
	delete(updates, "is_password_set")

	return repository.UpdateEmployeeFields(existingEmployee, updates)
}

func DeleteEmployee(employeeID int, companyID uint) error {
	// Verify employee belongs to this company before deleting
	existingEmployee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || existingEmployee == nil || existingEmployee.CompanyID != int(companyID) {
		return fmt.Errorf("employee not found or does not belong to your company")
	}

	return repository.DeleteEmployee(employeeID)
}

func GetPendingEmployeesByCompanyIDPaginated(companyID int, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
	return repository.GetPendingEmployeesByCompanyIDPaginated(companyID, search, page, pageSize)
}

func ResendPasswordEmail(employeeID int, companyID uint) error {
	// Verify employee exists and belongs to this company
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || employee == nil || employee.CompanyID != int(companyID) {
		return fmt.Errorf("employee not found or does not belong to your company")
	}

	// Check if employee has already set their password
	if employee.IsPasswordSet {
		return fmt.Errorf("employee has already set their password")
	}

	// Invalidate/delete any existing initial password tokens for this employee
	if err := repository.InvalidatePasswordResetTokensByUserID(uint(employee.ID), "employee_initial_password"); err != nil {
		log.Printf("Error invalidating old initial password tokens for employee %d: %v", employee.ID, err)
		// Continue, but log the error
	}

	// Generate new password reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24) // Token valid for 24 hours

	passwordResetToken := &models.PasswordResetTokenTable{
		UserID:    employee.ID,
		TokenType: "employee_initial_password",
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := repository.CreatePasswordResetToken(passwordResetToken); err != nil {
		log.Printf("Error creating new password reset token for employee %d: %v", employee.ID, err)
		return fmt.Errorf("failed to generate new initial password link")
	}

	// Send email with password reset link in a goroutine (background)
	resetLink := os.Getenv("FRONTEND_BASE_URL") + "/initial-password-setup?token=" + token
	if os.Getenv("FRONTEND_BASE_URL") == "" {
		resetLink = "http://localhost:5173/initial-password-setup?token=" + token // Fallback for development
	}

	go func(email, name, link string) {
		if err := helper.SendPasswordResetEmail(email, name, link); err != nil {
			log.Printf("Error sending initial password email to %s in background: %v", email, err)
		}
	}(employee.Email, employee.Name, resetLink)

	return nil
}



type BulkImportResult struct {
	RowNumber int    `json:"row_number"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func BulkCreateEmployees(ctx context.Context, companyID int, excelFile *excelize.File) ([]BulkImportResult, int, int, error) {
	rows, err := excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get rows from Excel sheet: %w", err)
	}

	if len(rows) <= 1 {
		return nil, 0, 0, fmt.Errorf("excel file is empty or only contains headers")
	}

	shifts, err := repository.GetShiftsByCompanyID(companyID)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching shifts for bulk import: %w", err)
	}
	shiftNameToID := make(map[string]int)
	for _, shift := range shifts {
		shiftNameToID[shift.Name] = shift.ID
	}

	results := []BulkImportResult{}
	successCount := 0
	failedCount := 0

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		rowNum := i + 1 // Human-readable row number

		// Ensure row has enough columns
		if len(row) < 5 {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Missing required columns."})
			failedCount++
			continue
		}

		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		position := strings.TrimSpace(row[2])
		employeeIDNumber := strings.TrimSpace(row[3])
		shiftName := strings.TrimSpace(row[4])

		// Basic validation
		if name == "" || email == "" || position == "" || employeeIDNumber == "" {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Name, Email, Position, or Employee ID Number cannot be empty."})
			failedCount++
			continue
		}

		// Validate email format (simple check)
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Invalid email format."})
			failedCount++
			continue
		}

		var shiftID *int
		if shiftName != "" {
			id, ok := shiftNameToID[shiftName]
				if !ok {
				results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: fmt.Sprintf("Shift name '%s' not found.", shiftName)})
				failedCount++
				continue
			}
			shiftID = &id
		} else {
			// If shift name is empty, try to assign default shift
			defaultShift, err := repository.GetDefaultShiftByCompanyID(companyID)
			if err != nil {
				log.Printf("No default shift found for company %d. Employee %s will be created without a shift.", companyID, email)
			} else if defaultShift != nil {
				shiftID = &defaultShift.ID
			}
		}

		// Check if employee already exists by email or employee_id_number
		existingEmployee, err := repository.GetEmployeeByEmailOrIDNumber(email, employeeIDNumber)
		if err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Database error during existence check."})
			failedCount++
			continue
		}
		if existingEmployee != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Employee with this email or ID number already exists."})
			failedCount++
			continue
		}

		// Check subscription limit before creating employee
		var company models.CompaniesTable
		if err := database.DB.Preload("SubscriptionPackage").Preload("CustomOffer").First(&company, uint(companyID)).Error; err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Failed to retrieve company subscription info."})
			failedCount++
			continue
		}

		var maxEmployeesLimit int
		if company.CustomOfferID != nil && company.CustomOffer != nil {
			maxEmployeesLimit = company.CustomOffer.MaxEmployees
		} else if company.SubscriptionPackage.ID != 0 {
			maxEmployeesLimit = company.SubscriptionPackage.MaxEmployees
		} else {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Company has no active subscription package or custom offer."})
			failedCount++
			continue
		}

		var currentEmployeeCount int64
		if err := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", companyID).Count(&currentEmployeeCount).Error; err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Failed to count existing employees for limit check."})
			failedCount++
			continue
		}

		if currentEmployeeCount >= int64(maxEmployeesLimit) {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Employee limit reached for your current plan."})
			failedCount++
			continue
		}

		// Create employee
		employee := &models.EmployeesTable{
			CompanyID:        companyID,
			Email:            email,
			Name:             name,
			Position:         position,
			EmployeeIDNumber: employeeIDNumber,
			ShiftID:          shiftID,
			Role:             "employee",
		}

		if err := repository.CreateEmployee(employee); err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Failed to create employee in database."})
			failedCount++
			continue
		}

		// Generate password reset token and send email in a goroutine
		go func(empID int, empEmail, empName string) {
			token := uuid.New().String()
			expiresAt := time.Now().Add(time.Hour * 24)
			passwordResetToken := &models.PasswordResetTokenTable{
				UserID:    empID,
				TokenType: "employee_initial_password",
				Token:     token,
				ExpiresAt: expiresAt,
			}
			if err := repository.CreatePasswordResetToken(passwordResetToken); err != nil {
				log.Printf("Error creating password reset token for employee %d: %v", empID, err)
				return
			}
			resetLink := os.Getenv("FRONTEND_BASE_URL") + "/initial-password-setup?token=" + token
			if os.Getenv("FRONTEND_BASE_URL") == "" {
				resetLink = "http://localhost:5173/initial-password-setup?token=" + token
			}
			if err := helper.SendPasswordResetEmail(empEmail, empName, resetLink); err != nil {
				log.Printf("Error sending initial password email to %s in background: %v", empEmail, err)
			}
		}(employee.ID, employee.Email, employee.Name)

		results = append(results, BulkImportResult{RowNumber: rowNum, Status: "success", Message: "Employee created successfully."})
		successCount++
	}

	return results, successCount, failedCount, nil
}

// PythonServerResponse defines the structure for the JSON response from the Python server.

type PythonServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendFaceRequestToPython handles the TCP communication with the Python face recognition server.
func SendFaceRequestToPython(action, base64ImageData, dbImagePath string) (*PythonServerResponse, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		return nil, fmt.Errorf("could not connect to python server: %w", err)
	}
	defer conn.Close()

	payload := map[string]string{
		"action":              action,
		"client_image_data": base64ImageData,
		"db_image_path":       dbImagePath,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json payload: %w", err)
	}

	_, err = conn.Write(append(requestBody, '\n'))
	if err != nil {
		return nil, fmt.Errorf("could not send data to python server: %w", err)
	}

	// Read the response from the server
	reader := bufio.NewReader(conn)
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("could not read response from python server: %w", err)
	}

	var response PythonServerResponse
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return nil, fmt.Errorf("could not unmarshal python server response: %w", err)
	}

	return &response, nil
}

func UploadFaceImage(employeeID int, companyID int, file *multipart.FileHeader) (string, error) {
	log.Printf("UploadFaceImage: Processing upload for EmployeeID: %d, CompanyID: %d", employeeID, companyID)

	// 1. Handle the image file from the form
	if file == nil {
		return "", fmt.Errorf("image file is required")
	}

	log.Printf("UploadFaceImage: Received file: %s, Size: %d", file.Filename, file.Size)

	// 2. Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		log.Printf("UploadFaceImage: Invalid file extension: %s", ext)
		return "", fmt.Errorf("invalid file type. Only JPG, JPEG, and PNG are allowed")
	}

	// 3. Read and encode the image for face detection
	openedFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded image: %w", err)
	}
	defer openedFile.Close()

	imageBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)

	// 4. Send to Python server for face validation
	faceCheckResult, err := SendFaceRequestToPython("check_face", encodedImage, "")
	if err != nil {
		log.Printf("[Go] Error communicating with Python server: %v", err)
		return "", fmt.Errorf("could not validate face: error communicating with recognition service")
	}

	// 5. Analyze the response from Python server
	if faceCheckResult.Status != "face_found" {
		log.Printf("[Go] Face check failed for employee %d: %s", employeeID, faceCheckResult.Message)
		return "", errors.New(faceCheckResult.Message) // Return Python's error message to the handler
	}

	log.Printf("[Go] Face check successful for employee %d: %s", employeeID, faceCheckResult.Message)

	// 6. Delete old face images if they exist (to ensure only one reference image)
	existingImages, err := repository.GetFaceImagesByEmployeeID(employeeID)
	if err != nil {
		log.Printf("UploadFaceImage: Could not check for existing images for employee %d: %v", employeeID, err)
		// Not a fatal error, so we continue
	}
	log.Printf("UploadFaceImage: Found %d existing images for employee %d", len(existingImages), employeeID)
	for _, img := range existingImages {
		if err := os.Remove(img.ImagePath); err != nil {
			log.Printf("UploadFaceImage: Failed to delete old image file %s: %v", img.ImagePath, err)
		}
		if err := repository.DeleteFaceImage(img.ID); err != nil {
			log.Printf("UploadFaceImage: Failed to delete old image record from DB %d: %v", img.ID, err)
		}
	}

	// 7. Save the new file using the helper function
	subDir := filepath.Join("employee_faces", strconv.Itoa(companyID))
	savePath, err := helper.SaveUploadedFile(file, subDir)
	if err != nil {
		return "", fmt.Errorf("failed to save face image file: %w", err)
	}
	log.Printf("UploadFaceImage: Saved new image to: %s", savePath)

	// 8. Record the new face image in the database
	faceImage := &models.FaceImagesTable{
		EmployeeID: employeeID,
		ImagePath:  savePath,
	}
	log.Printf("UploadFaceImage: Attempting to record face image in DB for EmployeeID: %d, ImagePath: %s", faceImage.EmployeeID, faceImage.ImagePath)
	if err := repository.CreateFaceImage(faceImage); err != nil {
		log.Printf("UploadFaceImage: Failed to record face image in database: %v", err)
		return "", fmt.Errorf("failed to record face image in database")
	}
	log.Printf("UploadFaceImage: Face image successfully recorded in DB with ID: %d", faceImage.ID)

	return savePath, nil
}

func GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	return repository.GetFaceImagesByEmployeeID(employeeID)
}

type UpdateEmployeeProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Position string `json:"position" binding:"required"`
}

func UpdateEmployeeProfile(employeeID int, req UpdateEmployeeProfileRequest) error {
	// Get existing employee to update
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || employee == nil {
		return fmt.Errorf("employee not found")
	}

	// Update fields
	employee.Name = req.Name
	employee.Email = req.Email
	employee.Position = req.Position

	// Save changes
	if err := repository.UpdateEmployee(employee); err != nil {
		return fmt.Errorf("failed to update employee profile: %w", err)
	}
	return nil
}

func ChangeEmployeePassword(employeeID int, oldPassword, newPassword, confirmNewPassword string) error {
	// Get existing employee
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || employee == nil {
		return fmt.Errorf("employee not found")
	}

	// Verify old password
	if helper.CheckPasswordHash(oldPassword, employee.Password) != nil {
		return fmt.Errorf("incorrect old password")
	}

	// Validate new password complexity
	if !helper.IsValidPassword(newPassword) {
		return fmt.Errorf("new password must be at least 8 characters long, contain uppercase, lowercase, and a number")
	}

	// Check if new password is the same as old password
	if newPassword == oldPassword {
		return fmt.Errorf("new password cannot be the same as the old password")
	}

	// Check if new password and confirmation match
	if newPassword != confirmNewPassword {
		return fmt.Errorf("new password and confirmation do not match")
	}

	// Update password
	if err := repository.UpdateEmployeePassword(employee, newPassword); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

type EmployeeDashboardSummary struct {
	EmployeeName            string                 `json:"employee_name"`
	EmployeePosition        string                 `json:"employee_position"`
	TodayAttendanceStatus   string                 `json:"today_attendance_status"`
	PendingLeaveRequestsCount int                    `json:"pending_leave_requests_count"`
	RecentAttendances       []models.AttendancesTable `json:"recent_attendances"`
}

func GetEmployeeDashboardSummary(employeeID int) (*EmployeeDashboardSummary, error) {
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil || employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	todayAttendance, err := repository.GetTodayAttendanceByEmployeeID(employeeID)
	var todayAttendanceStatus string
	if err != nil {
		log.Printf("Error getting today's attendance for employee %d: %v", employeeID, err)
		todayAttendanceStatus = "Unavailable"
	} else if todayAttendance != nil {
		todayAttendanceStatus = todayAttendance.Status
	} else {
		todayAttendanceStatus = "Not Checked In"
	}

	pendingLeaveRequests, err := repository.GetPendingLeaveRequestsByEmployeeID(employeeID)
	var pendingLeaveRequestsCount int
	if err != nil {
		log.Printf("Error getting pending leave requests for employee %d: %v", employeeID, err)
		pendingLeaveRequestsCount = 0
	} else {
		pendingLeaveRequestsCount = len(pendingLeaveRequests)
	}

	recentAttendances, err := repository.GetRecentAttendancesByEmployeeID(employeeID, 5)
	if err != nil {
		log.Printf("Error getting recent attendances for employee %d: %v", employeeID, err)
		recentAttendances = []models.AttendancesTable{}
	}

	response := &EmployeeDashboardSummary{
		EmployeeName:            employee.Name,
		EmployeePosition:        employee.Position,
		TodayAttendanceStatus:   todayAttendanceStatus,
		PendingLeaveRequestsCount: pendingLeaveRequestsCount,
		RecentAttendances:       recentAttendances,
	}

	return response, nil
}

// EmployeeProfileResponse defines the structure for the employee profile response.
type EmployeeProfileResponse struct {
	models.EmployeesTable
	Shift                      *models.ShiftsTable         `json:"shift,omitempty"`
	CompanyAttendanceLocations []models.AttendanceLocation `json:"company_attendance_locations,omitempty"`
	FaceImages                 []models.FaceImagesTable    `json:"face_images,omitempty"`
	FaceImageRegistered        bool                        `json:"face_image_registered"`
}

// GetEmployeeProfile handles fetching the profile for the currently logged-in employee.
func GetEmployeeProfile(employeeID int) (*EmployeeProfileResponse, error) {
	// Get employee data from repository
	employee, err := repository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee profile: %w", err)
	}
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	// Get shift data
	var shift *models.ShiftsTable
	if employee.ShiftID != nil {
		shift, _ = repository.GetShiftByID(*employee.ShiftID)
	}

	// Get company attendance locations
	locations, _ := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))

	// Get face images
	faceImages, _ := repository.GetFaceImagesByEmployeeID(employeeID)

	// Determine if face image is registered
	faceImageRegistered := len(faceImages) > 0

	// Create the response object
	profileResponse := &EmployeeProfileResponse{
		EmployeesTable:             *employee,
		Shift:                      shift,
		CompanyAttendanceLocations: locations,
		FaceImages:                 faceImages,
		FaceImageRegistered:        faceImageRegistered,
	}

	return profileResponse, nil
}


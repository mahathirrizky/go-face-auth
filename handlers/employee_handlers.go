package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- Employee Handlers ---

type CreateEmployeeRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Position       string `json:"position" binding:"required"`
	EmployeeIDNumber string `json:"employee_id_number" binding:"required"` // Added
	ShiftID        *int   `json:"shift_id"` // Optional: Pointer to int to allow null/omission
}

func CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

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

	// Retrieve company and its subscription package
	var company models.CompaniesTable
	if err := database.DB.Preload("SubscriptionPackage").First(&company, compID).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company information")
		return
	}

	// Check current employee count
	var employeeCount int64
	if err := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", compID).Count(&employeeCount).Error; err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to count existing employees")
		return
	}

	// Check if adding a new employee would exceed the package limit
	if employeeCount >= int64(company.SubscriptionPackage.MaxEmployees) {
		helper.SendError(c, http.StatusForbidden, "Employee limit reached for your subscription package")
		return
	}

	employee := &models.EmployeesTable{
		CompanyID: compID,
		Email:     req.Email,
		Name:      req.Name,
		Position:  req.Position,
		EmployeeIDNumber: req.EmployeeIDNumber, // Added
		Role:      "employee", // Set default role to employee
	}

	// Determine the shift ID for the employee
	if req.ShiftID != nil {
		employee.ShiftID = req.ShiftID
	} else {
		// If no shift ID is provided, try to find the default shift for the company
		defaultShift, err := repository.GetDefaultShiftByCompanyID(compID)
		if err != nil {
			// If no default shift is found, log and continue without assigning a shift
			log.Printf("No default shift found for company %d. Employee %s will be created without a shift.", compID, req.Email)
		} else if defaultShift != nil {
			employee.ShiftID = &defaultShift.ID
		}
	}

	if err := repository.CreateEmployee(employee); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create employee")
		return
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
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate initial password link.")
		return
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
	}(employee.Email, employee.Name, resetLink) // Pass values to goroutine to avoid closure issues

	helper.SendSuccess(c, http.StatusCreated, "Employee created successfully. An email with initial password setup link has been sent.", gin.H{"employee_id": employee.ID, "employee_email": employee.Email})
}

func GetEmployeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("employeeID"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	employee, err := repository.GetEmployeeByID(id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee.")
		return
	}

	if employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee retrieved successfully.", employee)
}

func GetEmployeesByCompanyID(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	employees, err := repository.GetEmployeesByCompanyID(companyID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employees.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employees retrieved successfully.", employees)
}

// SearchEmployees handles searching for employees by name within a specific company.
func SearchEmployees(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	name := c.Query("name")

	employees, err := repository.SearchEmployees(companyID, name)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to search employees")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employees found successfully", employees)
}

// UpdateEmployee handles updating an existing employee.
func UpdateEmployee(c *gin.Context) {
	idStr := c.Param("employeeID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	var employee models.EmployeesTable
	if err := c.ShouldBindJSON(&employee); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	employee.ID = id // Ensure the ID from the URL is used

	// Get company ID from JWT claims to ensure employee belongs to the company
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

	// Verify employee belongs to this company
	existingEmployee, err := repository.GetEmployeeByID(id)
	if err != nil || existingEmployee == nil || existingEmployee.CompanyID != compID {
		helper.SendError(c, http.StatusNotFound, "Employee not found or does not belong to your company.")
		return
	}

	// Preserve password if not provided in update request
	if employee.Password == "" {
		employee.Password = existingEmployee.Password
	}

	if err := repository.UpdateEmployee(&employee); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update employee.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee updated successfully.", employee)
}

// DeleteEmployee handles deleting an employee.
func DeleteEmployee(c *gin.Context) {
	idStr := c.Param("employeeID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	// Get company ID from JWT claims to ensure employee belongs to the company
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

	// Verify employee belongs to this company before deleting
	existingEmployee, err := repository.GetEmployeeByID(id)
	if err != nil || existingEmployee == nil || existingEmployee.CompanyID != compID {
		helper.SendError(c, http.StatusNotFound, "Employee not found or does not belong to your company.")
		return
	}

	if err := repository.DeleteEmployee(id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete employee.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee deleted successfully.", nil)
}

// GetPendingEmployees handles fetching employees who have not set their password yet.
func GetPendingEmployees(c *gin.Context) {
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

	employees, err := repository.GetPendingEmployees(compID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve pending employees.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Pending employees retrieved successfully.", employees)
}

// ResendPasswordEmailRequest defines the structure for a resend password email request.
func ResendPasswordEmail(c *gin.Context) {
	idStr := c.Param("employee_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

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

	// Verify employee exists and belongs to this company
	employee, err := repository.GetEmployeeByID(id)
	if err != nil || employee == nil || employee.CompanyID != compID {
		helper.SendError(c, http.StatusNotFound, "Employee not found or does not belong to your company.")
		return
	}

	// Check if employee has already set their password
	if employee.IsPasswordSet {
		helper.SendError(c, http.StatusBadRequest, "Employee has already set their password.")
		return
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
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate new initial password link.")
		return
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

	helper.SendSuccess(c, http.StatusOK, "Initial password setup email resent successfully.", nil)
}

// --- Face Image Handlers ---

// UploadFaceImage handles the initial upload of a face image for the authenticated employee.
func UploadFaceImage(c *gin.Context) {
	// 1. Get Employee and Company ID from JWT Token (Security Best Practice)
	employeeIDFromToken, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token.")
		return
	}
	empIDFloat, _ := employeeIDFromToken.(float64)
	empID := int(empIDFloat)

	companyIDFromToken, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token.")
		return
	}
	compIDFloat, _ := companyIDFromToken.(float64)
	compID := int(compIDFloat)

	// 2. Handle the image file from the form
	file, err := c.FormFile("image")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Image file is required.")
		return
	}

	// 3. Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		helper.SendError(c, http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed.")
		return
	}

	// 4. Delete old face images if they exist (to ensure only one reference image)
	existingImages, err := repository.GetFaceImagesByEmployeeID(empID)
	if err != nil {
		log.Printf("Could not check for existing images for employee %d: %v", empID, err)
		// Not a fatal error, so we continue
	}
	for _, img := range existingImages {
		if err := os.Remove(img.ImagePath); err != nil {
			log.Printf("Failed to delete old image file %s: %v", img.ImagePath, err)
		}
		if err := repository.DeleteFaceImage(img.ID); err != nil {
			log.Printf("Failed to delete old image record from DB %d: %v", img.ID, err)
		}
	}

	// 5. Create a unique filename and path
	companyDir := filepath.Join("images", "employee_faces", strconv.Itoa(compID))
	if err := os.MkdirAll(companyDir, os.ModePerm); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create image directory.")
		return
	}
	uniqueFilename := uuid.New().String() + ext
	savePath := filepath.Join(companyDir, uniqueFilename)

	// 6. Save the new file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to save image file.")
		return
	}

	// 7. Record the new face image in the database
	faceImage := &models.FaceImagesTable{
		EmployeeID: empID,
		ImagePath:  savePath,
	}
	if err := repository.CreateFaceImage(faceImage); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record face image in database.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Face image uploaded successfully.", gin.H{
		"employee_id": empID,
		"image_path":  savePath,
	})
}

func GetFaceImagesByEmployeeID(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("employee_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	faceImages, err := repository.GetFaceImagesByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve face images.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Face images retrieved successfully.", faceImages)
}

// GetEmployeeDashboardSummary handles fetching the dashboard summary for the logged-in employee.
func GetEmployeeDashboardSummary(c *gin.Context) {
	// Get employee ID from JWT claims
	employeeIDFromContext, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token claims.")
		return
	}
	empIDFloat, ok := employeeIDFromContext.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := int(empIDFloat)

	// Get employee data
	employee, err := repository.GetEmployeeByID(empID)
	if err != nil || employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get today's attendance status
	todayAttendance, err := repository.GetTodayAttendanceByEmployeeID(empID)
	var todayAttendanceStatus string
	if err != nil {
		log.Printf("Error getting today's attendance for employee %d: %v", empID, err)
		todayAttendanceStatus = "Tidak tersedia"
	} else if todayAttendance != nil {
		todayAttendanceStatus = todayAttendance.Status
	} else {
		todayAttendanceStatus = "Belum Absen"
	}

	// Get pending leave requests count
	pendingLeaveRequests, err := repository.GetPendingLeaveRequestsByEmployeeID(empID)
	var pendingLeaveRequestsCount int
	if err != nil {
		log.Printf("Error getting pending leave requests for employee %d: %v", empID, err)
		pendingLeaveRequestsCount = 0
	} else {
		pendingLeaveRequestsCount = len(pendingLeaveRequests)
	}

	// Get recent attendances (e.g., last 5)
	recentAttendances, err := repository.GetRecentAttendancesByEmployeeID(empID, 5)
	if err != nil {
		log.Printf("Error getting recent attendances for employee %d: %v", empID, err)
		recentAttendances = []models.AttendancesTable{}
	}

	// Prepare response data
	response := gin.H{
		"employee_name":             employee.Name,
		"employee_position":         employee.Position,
		"today_attendance_status":   todayAttendanceStatus,
		"pending_leave_requests_count": pendingLeaveRequestsCount,
		"recent_attendances":        recentAttendances,
	}

	helper.SendSuccess(c, http.StatusOK, "Employee dashboard summary retrieved successfully.", response)
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
func GetEmployeeProfile(c *gin.Context) {
	// Get employee ID from JWT claims
	employeeIDFromContext, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token claims.")
		return
	}
	empIDFloat, ok := employeeIDFromContext.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := int(empIDFloat)

	// Get employee data from repository
	employee, err := repository.GetEmployeeByID(empID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employee profile.")
		return
	}
	if employee == nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	// Get shift data
	var shift *models.ShiftsTable
	if employee.ShiftID != nil {
		shift, _ = repository.GetShiftByID(*employee.ShiftID)
	}

	// Get company attendance locations
	locations, _ := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))

	// Get face images
	faceImages, _ := repository.GetFaceImagesByEmployeeID(empID)

	// Determine if face image is registered
	faceImageRegistered := len(faceImages) > 0

	// Create the response object
	profileResponse := EmployeeProfileResponse{
		EmployeesTable:             *employee,
		Shift:                      shift,
		CompanyAttendanceLocations: locations,
		FaceImages:                 faceImages,
		FaceImageRegistered:        faceImageRegistered,
	}

	helper.SendSuccess(c, http.StatusOK, "Profile retrieved successfully.", profileResponse)
}

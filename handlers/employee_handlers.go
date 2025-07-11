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
	resetLink := os.Getenv("FRONTEND_BASE_URL") + "/reset-password?token=" + token
	if os.Getenv("FRONTEND_BASE_URL") == "" {
		resetLink = "http://localhost:5173/reset-password?token=" + token // Fallback for development
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
	id, err := strconv.Atoi(c.Param("id"))
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
	idStr := c.Param("id")
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
	idStr := c.Param("id")
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
type ResendPasswordEmailRequest struct {
	EmployeeID uint `json:"employee_id" binding:"required"`
}

// ResendPasswordEmail handles resending the initial password setup email to an employee.
func ResendPasswordEmail(c *gin.Context) {
	var req ResendPasswordEmailRequest
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

	// Verify employee exists and belongs to this company
	employee, err := repository.GetEmployeeByID(int(req.EmployeeID))
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
	resetLink := os.Getenv("FRONTEND_BASE_URL") + "/reset-password?token=" + token
	if os.Getenv("FRONTEND_BASE_URL") == "" {
		resetLink = "http://localhost:5173/reset-password?token=" + token // Fallback for development
	}

	go func(email, name, link string) {
		if err := helper.SendPasswordResetEmail(email, name, link); err != nil {
			log.Printf("Error sending initial password email to %s in background: %v", email, err)
		}
	}(employee.Email, employee.Name, resetLink)

	helper.SendSuccess(c, http.StatusOK, "Initial password setup email resent successfully.", nil)
}

// --- Face Image Handlers ---

func UploadFaceImage(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.PostForm("employee_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

	// Check if employee exists
	_, err = repository.GetEmployeeByID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, "Employee not found.")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Image file is required.")
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		helper.SendError(c, http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed.")
		return
	}

	// Create a unique filename
	filename := "employee_" + strconv.Itoa(employeeID) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
	
	// Define the path to save the image
	savePath := filepath.Join("images", "employee_faces", filename)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create image directory.")
		return
	}

	// Save the file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to save image file.")
		return
	}

	faceImage := &models.FaceImagesTable{
		EmployeeID: employeeID,
		ImagePath:  savePath,
	}

	if err := repository.CreateFaceImage(faceImage); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to record face image in database.")
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Face image uploaded and recorded successfully.", gin.H{
		"employee_id": employeeID,
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

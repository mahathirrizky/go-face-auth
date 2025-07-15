package handlers

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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

// GenerateEmployeeTemplate generates an Excel template for bulk employee import.
func GenerateEmployeeTemplate(c *gin.Context) {
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

    // Fetch shifts for the company
    shifts, err := repository.GetShiftsByCompanyID(compID)
    if err != nil {
        log.Printf("Error fetching shifts for template generation: %v", err)
        helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve shifts for template.")
        return
    }

    // Debug: Print shifts
    log.Printf("Shifts retrieved: %+v", shifts)
    if len(shifts) == 0 {
        log.Printf("Warning: No shifts found for company ID %d", compID)
    }

    f := excelize.NewFile()
    // Create a sheet for shift names (not hidden for debugging)
    shiftSheetName := "ShiftData"
    f.NewSheet(shiftSheetName)
    // Tidak menyembunyikan sheet ShiftData untuk mempermudah pemeriksaan
    // f.SetSheetVisible(shiftSheetName, false)

    // Populate shift names in the ShiftData sheet and log them
    for i, shift := range shifts {
        shiftName := strings.TrimSpace(shift.Name) // Clean the shift name
        if shiftName == "" {
            log.Printf("Warning: Empty shift name at index %d", i)
            continue
        }
        cell := fmt.Sprintf("A%d", i+1)
        f.SetCellValue(shiftSheetName, cell, shiftName)
        log.Printf("Writing shift to ShiftData sheet: %s at cell %s", shiftName, cell)
    }

    // Debug: Verify ShiftData sheet content
    rows, err := f.GetRows(shiftSheetName)
    if err != nil {
        log.Printf("Error reading ShiftData sheet: %v", err)
    } else {
        log.Printf("ShiftData sheet content: %+v", rows)
    }

    // Set the main sheet
    mainSheetName := "Employees"
    f.SetSheetName("Sheet1", mainSheetName)

    // Set headers for the main sheet
    headers := []string{"Name", "Email", "Position", "Employee ID Number", "Shift Name"}
    for i, header := range headers {
        f.SetCellValue(mainSheetName, fmt.Sprintf("%s1", string(rune('A'+i))), header)
    }

    // Apply data validation for Shift Name column (e.g., for first 100 rows)
    if len(shifts) > 0 {
        // Create a list of shift names for the dropdown
        shiftNames := make([]string, 0, len(shifts))
        for _, shift := range shifts {
            if name := strings.TrimSpace(shift.Name); name != "" {
                shiftNames = append(shiftNames, name)
            }
        }
        log.Printf("Shift names for dropdown: %+v", shiftNames)

        // Create data validation for dropdown
        dv := excelize.NewDataValidation(true)
        dv.Sqref = "E2:E101" // Apply to Shift Name column (E) from row 2 to 101
        dv.SetDropList(shiftNames) // Use SetDropList for direct list
        dv.ShowDropDown = true
        dv.AllowBlank = true


        if err := f.AddDataValidation(mainSheetName, dv); err != nil {
            log.Printf("Error setting data validation: %v", err)
            helper.SendError(c, http.StatusInternalServerError, "Failed to set data validation in Excel.")
            return
        }
        log.Printf("Data validation applied successfully to %s", dv.Sqref)
    } else {
        log.Printf("No data validation applied because no shifts are available")
    }

    // Save the file locally for debugging
    if err := f.SaveAs("debug_employee_template.xlsx"); err != nil {
        log.Printf("Error saving debug Excel file: %v", err)
    } else {
        log.Printf("Debug Excel file saved as debug_employee_template.xlsx")
    }

    // Set response headers for Excel file download
    c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    c.Header("Content-Disposition", "attachment; filename=employee_template.xlsx")

    // Write the Excel file to the response writer
    if err := f.Write(c.Writer); err != nil {
        log.Printf("Error writing excel file: %v", err)
        helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
        return
    }
    log.Printf("Excel file generated successfully")
}
// BulkCreateEmployees handles bulk creation of employees from an uploaded Excel file.
func BulkCreateEmployees(c *gin.Context) {
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

	file, err := c.FormFile("file")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "No file uploaded.")
		return
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to open uploaded file.")
		return
	}
	defer f.Close()

	// Read the Excel file
	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Failed to read Excel file: "+err.Error())
		return
	}

	// Get all rows from the first sheet
	rows, err := excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to get rows from Excel sheet.")
		return
	}

	if len(rows) <= 1 {
		helper.SendError(c, http.StatusBadRequest, "Excel file is empty or only contains headers.")
		return
	}

	// Fetch shifts for the company to map shift names to IDs
	shifts, err := repository.GetShiftsByCompanyID(compID)
	if err != nil {
		log.Printf("Error fetching shifts for bulk import: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve shifts for processing.")
		return
	}
	shiftNameToID := make(map[string]int)
	for _, shift := range shifts {
		shiftNameToID[shift.Name] = shift.ID
	}

	// Prepare results
	type BulkImportResult struct {
		RowNumber int    `json:"row_number"`
		Status    string `json:"status"`
		Message   string `json:"message"`
	}
	results := []BulkImportResult{}
	successCount := 0
	failedCount := 0

	// Process rows (skip header row)
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
			defaultShift, err := repository.GetDefaultShiftByCompanyID(compID)
			if err != nil {
				log.Printf("No default shift found for company %d. Employee %s will be created without a shift.", compID, email)
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
		if err := database.DB.Preload("SubscriptionPackage").First(&company, compID).Error; err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Failed to retrieve company subscription info."})
			failedCount++
			continue
		}

		var currentEmployeeCount int64
		if err := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", compID).Count(&currentEmployeeCount).Error; err != nil {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Failed to count existing employees for limit check."})
			failedCount++
			continue
		}

		if currentEmployeeCount >= int64(company.SubscriptionPackage.MaxEmployees) {
			results = append(results, BulkImportResult{RowNumber: rowNum, Status: "failed", Message: "Employee limit reached for your subscription package."})
			failedCount++
			continue
		}

		// Create employee
		employee := &models.EmployeesTable{
			CompanyID:        compID,
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

	helper.SendSuccess(c, http.StatusOK, "Bulk import complete.", gin.H{
		"total_processed": len(rows) - 1,
		"success_count":   successCount,
		"failed_count":    failedCount,
		"results":         results,
	})
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

	log.Printf("UploadFaceImage: Processing upload for EmployeeID: %d, CompanyID: %d", empID, compID)

	// 2. Handle the image file from the form
	file, err := c.FormFile("face_image") // Changed from "image" to "face_image"
	if err != nil {
		log.Printf("UploadFaceImage: Error getting form file: %v", err)
		helper.SendError(c, http.StatusBadRequest, "Image file is required.")
		return
	}

	log.Printf("UploadFaceImage: Received file: %s, Size: %d", file.Filename, file.Size)

	// 3. Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		log.Printf("UploadFaceImage: Invalid file extension: %s", ext)
		helper.SendError(c, http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed.")
		return
	}

	// 4. Delete old face images if they exist (to ensure only one reference image)
	existingImages, err := repository.GetFaceImagesByEmployeeID(empID)
	if err != nil {
		log.Printf("UploadFaceImage: Could not check for existing images for employee %d: %v", empID, err)
		// Not a fatal error, so we continue
	}
	log.Printf("UploadFaceImage: Found %d existing images for employee %d", len(existingImages), empID)
	for _, img := range existingImages {
		if err := os.Remove(img.ImagePath); err != nil {
			log.Printf("UploadFaceImage: Failed to delete old image file %s: %v", img.ImagePath, err)
		}
		if err := repository.DeleteFaceImage(img.ID); err != nil {
			log.Printf("UploadFaceImage: Failed to delete old image record from DB %d: %v", img.ID, err)
		}
	}

	// 5. Create a unique filename and path
	storageBaseDir := os.Getenv("STORAGE_BASE_PATH")
	if storageBaseDir == "" {
		storageBaseDir = "/tmp/go_face_auth_data" // Fallback for development/testing
	}
	companyDir := filepath.Join(storageBaseDir, "employee_faces", strconv.Itoa(compID))
	if err := os.MkdirAll(companyDir, os.ModePerm); err != nil {
		log.Printf("UploadFaceImage: Failed to create image directory %s: %v", companyDir, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to create image directory.")
		return
	}
	uniqueFilename := uuid.New().String() + ext
	savePath := filepath.Join(companyDir, uniqueFilename)
	log.Printf("UploadFaceImage: Saving new image to: %s", savePath)

	// 6. Save the new file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		log.Printf("UploadFaceImage: Failed to save image file %s: %v", savePath, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to save image file.")
		return
	}

	// 7. Record the new face image in the database
	faceImage := &models.FaceImagesTable{
		EmployeeID: empID,
		ImagePath:  savePath,
	}
	log.Printf("UploadFaceImage: Attempting to record face image in DB for EmployeeID: %d, ImagePath: %s", faceImage.EmployeeID, faceImage.ImagePath)
	if err := repository.CreateFaceImage(faceImage); err != nil {
		log.Printf("UploadFaceImage: Failed to record face image in database: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to record face image in database.")
		return
	}
	log.Printf("UploadFaceImage: Face image successfully recorded in DB with ID: %d", faceImage.ID)

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

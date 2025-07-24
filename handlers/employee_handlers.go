package handlers

import (
	"go-face-auth/helper"
	"go-face-auth/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// --- Employee Handlers ---

func CreateEmployee(c *gin.Context) {
	var req services.CreateEmployeeRequest
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
	compID := uint(compIDFloat)

	employee, err := services.CreateEmployee(c.Request.Context(), compID, req)
	if err != nil {
		// Check for specific error messages from the service
		if err.Error() == "employee limit reached for your subscription package" {
			helper.SendError(c, http.StatusForbidden, err.Error())
		} else {
			helper.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Employee created successfully. An email with initial password setup link has been sent.", gin.H{"employee_id": employee.ID, "employee_email": employee.Email})
}

func GetEmployeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("employeeID"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid employee ID.")
		return
	}

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
	compID := uint(compIDFloat)

	employee, err := services.GetEmployeeByID(id, compID)
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	employees, totalRecords, err := services.GetEmployeesByCompanyIDPaginated(companyID, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve employees.")
		return
	}

	paginatedData := gin.H{
		"items":         employees,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Employees retrieved successfully.", paginatedData)
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

	employees, err := services.SearchEmployees(companyID, name)
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

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body.")
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
	compID := uint(compIDFloat)

	if err := services.UpdateEmployee(id, compID, updates); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update employee: "+err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee updated successfully.", nil)
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
	compID := uint(compIDFloat)

	if err := services.DeleteEmployee(id, compID); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to delete employee.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee deleted successfully.", nil)
}

// GetPendingEmployees handles fetching employees who have not set their password yet.
func GetPendingEmployees(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid company ID.")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	employees, totalRecords, err := services.GetPendingEmployeesByCompanyIDPaginated(companyID, search, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve pending employees.")
		return
	}

	paginatedData := gin.H{
		"items":         employees,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Pending employees retrieved successfully.", paginatedData)
}

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
	compID := uint(compIDFloat)

	if err := services.ResendPasswordEmail(id, compID); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

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
    shifts, err := services.GetShiftsByCompanyID(compID)
    if err != nil {
        log.Printf("Error fetching shifts for template generation: %v", err)
        helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve shifts for template.")
        return
    }

    f := excelize.NewFile()
    mainSheetName := "Employees"
    f.SetSheetName("Sheet1", mainSheetName)

    // Set headers for the main sheet
    headers := []string{"Name", "Email", "Position", "Employee ID Number", "Shift Name"}
    for i, header := range headers {
        cell, _ := excelize.CoordinatesToCellName(i+1, 1)
        f.SetCellValue(mainSheetName, cell, header)
    }

    if len(shifts) > 0 {
        shiftNames := make([]string, len(shifts))
        for i, s := range shifts {
            shiftNames[i] = s.Name
        }

        dv := excelize.NewDataValidation(true)
        dv.SetSqref("E2:E101")
        dv.SetDropList(shiftNames)

        if err := f.AddDataValidation(mainSheetName, dv); err != nil {
            log.Printf("Error adding data validation to sheet: %v", err)
            helper.SendError(c, http.StatusInternalServerError, "Failed to apply validation to Excel template.")
            return
        }
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

	results, successCount, failedCount, err := services.BulkCreateEmployees(c.Request.Context(), compID, excelFile)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Bulk import complete.", gin.H{
		"total_processed": successCount + failedCount,
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

	file, err := c.FormFile("face_image") // Changed from "image" to "face_image"
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, "Image file is required.")
		return
	}

	savePath, err := services.UploadFaceImage(empID, compID, file)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
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

	faceImages, err := services.GetFaceImagesByEmployeeID(employeeID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve face images.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Face images retrieved successfully.", faceImages)
}

// UpdateEmployeeProfile handles updating the profile of the currently logged-in employee.
func UpdateEmployeeProfile(c *gin.Context) {
	var req services.UpdateEmployeeProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

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

	if err := services.UpdateEmployeeProfile(empID, req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee profile updated successfully.", nil)
}

// ChangeEmployeePassword handles changing the password for the currently logged-in employee.
func ChangeEmployeePassword(c *gin.Context) {
	var req struct {
		OldPassword         string `json:"old_password" binding:"required"`
		NewPassword         string `json:"new_password" binding:"required"`
		ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

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

	if err := services.ChangeEmployeePassword(empID, req.OldPassword, req.NewPassword, req.ConfirmNewPassword); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Password changed successfully.", nil)
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

	summary, err := services.GetEmployeeDashboardSummary(empID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Employee dashboard summary retrieved successfully.", summary)
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

	profileResponse, err := services.GetEmployeeProfile(empID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Profile retrieved successfully.", profileResponse)
}
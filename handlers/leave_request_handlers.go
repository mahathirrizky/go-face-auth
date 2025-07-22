package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-face-auth/helper"

	"go-face-auth/services"
	"go-face-auth/websocket"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type CreateLeaveRequestPayload struct {
	Type      string `json:"type" binding:"required,oneof=cuti sakit"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02"`
	Reason    string `json:"reason" binding:"required,min=10"`
}

// Employee Handlers

func ApplyLeave(c *gin.Context) {
	employeeID, exists := c.Get("employeeID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token")
		return
	}
	empIDFloat, ok := employeeID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := uint(empIDFloat)

	var req CreateLeaveRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
		return
	}

	leaveRequest, err := services.ApplyLeave(empID, req.Type, req.StartDate, req.EndDate, req.Reason)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Leave request submitted successfully.", leaveRequest)
}

func GetMyLeaveRequests(c *gin.Context) {
	employeeID, exists := c.Get("id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Employee ID not found in token")
		return
	}
	empIDFloat, ok := employeeID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid employee ID type in token claims.")
		return
	}
	empID := uint(empIDFloat)

	var startDate *time.Time
	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD.")
			return
		}
		startDate = &parsedDate
	}

	var endDate *time.Time
	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD.")
			return
		}
		endDate = &parsedDate
	}

	leaveRequests, err := services.GetMyLeaveRequests(empID, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve leave requests.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Leave requests retrieved successfully.", leaveRequests)
}

// Admin Handlers

func GetAllCompanyLeaveRequests(c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	compIDFloat, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}
	compID := int(compIDFloat)

	// Get query params for pagination and filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
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

	leaveRequests, totalRecords, err := services.GetAllCompanyLeaveRequests(compID, status, search, startDate, endDate, page, pageSize)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve leave requests.")
		return
	}

	paginatedData := gin.H{
		"items":         leaveRequests,
		"total_records": totalRecords,
	}

	helper.SendSuccess(c, http.StatusOK, "Leave requests retrieved successfully.", paginatedData)
}

type ReviewLeaveRequestPayload struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

func ReviewLeaveRequest(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		leaveRequestID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, "Invalid leave request ID.")
			return
		}

		adminID, exists := c.Get("id") // Assuming adminID is set in JWT for admin users
		if !exists {
			helper.SendError(c, http.StatusUnauthorized, "Admin ID not found in token.")
			return
		}
		adminIDUint, ok := adminID.(float64)
		if !ok {
			helper.SendError(c, http.StatusInternalServerError, "Invalid admin ID type in token claims.")
			return
		}
		adminIDVal := uint(adminIDUint)

		var req ReviewLeaveRequestPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.SendError(c, http.StatusBadRequest, helper.GetValidationError(err))
			return
		}

		leaveRequest, err := services.ReviewLeaveRequest(uint(leaveRequestID), adminIDVal, req.Status)
		if err != nil {
			helper.SendError(c, http.StatusForbidden, err.Error())
			return
		}

		// Get company ID from admin's token
		companyID, exists := c.Get("companyID")
		if !exists {
			return // Should not happen if AuthMiddleware is used
		}
		compIDFloat, ok := companyID.(float64)
		if !ok {
			return // Should not happen
		}
		compID := int(compIDFloat)

		// Trigger dashboard update for the company
		go func() {
			summary, err := services.GetDashboardSummaryData(compID)
			if err != nil {
				log.Printf("Error fetching dashboard summary for WebSocket update after leave review: %v", err)
				return
			}
			hub.SendDashboardUpdate(compID, summary)
		}()

		helper.SendSuccess(c, http.StatusOK, "Leave request status updated successfully.", leaveRequest)
	}
}

// ExportCompanyLeaveRequestsToExcel exports all leave request records for the company to an Excel file.
func ExportCompanyLeaveRequestsToExcel(c *gin.Context) {
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

	status := c.Query("status")
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

	leaveRequests, err := services.ExportCompanyLeaveRequests(compID, status, search, startDate, endDate)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve leave requests for export.")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	sheetName := "Leave Requests"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers
	f.SetCellValue(sheetName, "A1", "Employee Name")
	f.SetCellValue(sheetName, "B1", "Type")
	f.SetCellValue(sheetName, "C1", "Start Date")
	f.SetCellValue(sheetName, "D1", "End Date")
	f.SetCellValue(sheetName, "E1", "Reason")
	f.SetCellValue(sheetName, "F1", "Status")

	// Apply style to header row
	style, err := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#DDEBF7"}}, // Light blue background
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		log.Printf("Error creating style: %v", err)
	} else {
		f.SetCellStyle(sheetName, "A1", "F1", style)
	}

	// Populate data
	for i, lr := range leaveRequests {
		row := i + 2 // Start from row 2 after headers
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), lr.Employee.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), lr.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), lr.StartDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), lr.EndDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), lr.Reason)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), lr.Status)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	fileName := "company_leave_requests.xlsx"
	dateRange := ""
	if startDate != nil && endDate != nil {
		dateRange = fmt.Sprintf("_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		dateRange = fmt.Sprintf("_%s_onwards", startDate.Format("2006-01-02"))
	} else if endDate != nil {
		dateRange = fmt.Sprintf("_until_%s", endDate.Format("2006-01-02"))
	}
	fileName = fmt.Sprintf("company_leave_requests%s.xlsx", dateRange)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel file to response: %v", err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to generate Excel file.")
		return
	}
}

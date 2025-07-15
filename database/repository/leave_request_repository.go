package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CreateLeaveRequest inserts a new leave request into the database.
func CreateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := database.DB.Create(leaveRequest)
	if result.Error != nil {
		log.Printf("Error creating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request created with ID: %d", leaveRequest.ID)
	return nil
}

// GetLeaveRequestByID retrieves a leave request by its ID.
func GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	var leaveRequest models.LeaveRequest
	result := database.DB.Preload("Employee").First(&leaveRequest, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Leave request not found
		}
		log.Printf("Error getting leave request with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &leaveRequest, nil
}

// GetAllLeaveRequests retrieves all leave requests.
func GetAllLeaveRequests() ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := database.DB.Preload("Employee").Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting all leave requests: %v", result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// GetLeaveRequestsByEmployeeID retrieves all leave requests for a given employee ID, optionally filtered by date range.
func GetLeaveRequestsByEmployeeID(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	query := database.DB.Preload("Employee").Where("employee_id = ?", employeeID)

	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", endDate)
	}

	result := query.Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting leave requests for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// UpdateLeaveRequest updates an existing leave request record in the database.
func UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := database.DB.Save(leaveRequest)
	if result.Error != nil {
		log.Printf("Error updating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request updated with ID: %d", leaveRequest.ID)
	return nil
}

// GetRecentLeaveRequestsByCompanyID retrieves recent leave requests for a given company ID.
func GetRecentLeaveRequestsByCompanyID(companyID int, limit int) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	log.Printf("Repository: Fetching recent leave requests for company %d, limit %d", companyID, limit)
	result := database.DB.Preload("Employee").Joins("join employees_tables on leave_requests.employee_id = employees_tables.id").Where("employees_tables.company_id = ?", companyID).Order("created_at DESC").Limit(limit).Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Repository: Error getting recent leave requests for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	log.Printf("Repository: Found %d recent leave requests for company %d.", len(leaveRequests), companyID)
	for i, lr := range leaveRequests {
		log.Printf("Repository: LeaveRequest %d - EmployeeID: %d, EmployeeName: %s", i, lr.EmployeeID, lr.Employee.Name)
	}
	return leaveRequests, nil
}

// IsEmployeeOnApprovedLeave checks if an employee has an approved leave or sick request for a specific date.
func IsEmployeeOnApprovedLeave(employeeID int, date time.Time) (bool, error) {
	var count int64
	// Normalize date to start of day for comparison
	checkDate := date.Truncate(24 * time.Hour)

	result := database.DB.Model(&models.LeaveRequest{}).Where("employee_id = ? AND status = ? AND start_date <= ? AND end_date >= ?", employeeID, "approved", checkDate, checkDate).Count(&count)
	if result.Error != nil {
		log.Printf("Error checking approved leave for employee %d on date %s: %v", employeeID, date.Format("2006-01-02"), result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

// IsEmployeeOnApprovedLeaveDateRange checks if an employee has an approved leave or sick request within a specific date range.
func IsEmployeeOnApprovedLeaveDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	var count int64
	query := database.DB.Model(&models.LeaveRequest{}).Where("employee_id = ? AND status = ?", employeeID, "approved")

	// Check for overlap with the requested date range
	if startDate != nil {
		query = query.Where("end_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("start_date <= ?", *endDate)
	}

	result := query.Count(&count)
	if result.Error != nil {
		log.Printf("Error checking approved leave for employee %d in date range: %v", employeeID, result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

// GetPendingLeaveRequestsByEmployeeID retrieves all pending leave requests for a given employee ID.
func GetPendingLeaveRequestsByEmployeeID(employeeID int) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := database.DB.Preload("Employee").Where("employee_id = ? AND status = ?", employeeID, "pending").Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting pending leave requests for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

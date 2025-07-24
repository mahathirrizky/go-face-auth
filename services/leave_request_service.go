package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"time"
	"mime/multipart"
	"go-face-auth/helper"
	"path/filepath"
	"strconv"
)

func ApplyLeave(employeeID uint, leaveType, startDateStr, endDateStr, reason string, sickNote *multipart.FileHeader) (*models.LeaveRequest, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format. Use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format. Use YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	employee, err := repository.GetEmployeeByID(int(employeeID))
	if err != nil || employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	var sickNotePath string
	if leaveType == "sakit" && sickNote != nil {
		// Save the sick note file
		subDir := filepath.Join("sick_notes", strconv.Itoa(employee.CompanyID), strconv.Itoa(int(employeeID)))
		sickNotePath, err = helper.SaveUploadedFile(sickNote, subDir)
		if err != nil {
			return nil, fmt.Errorf("failed to save sick note file: %w", err)
		}
	}

	leaveRequest := &models.LeaveRequest{
		EmployeeID: employeeID,
		Type:       leaveType,
		StartDate:  startDate,
		EndDate:    endDate,
		Reason:     reason,
		Status:     "pending",
		SickNotePath: sickNotePath,
	}

	if err := repository.CreateLeaveRequest(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to submit leave request: %w", err)
	}

	return leaveRequest, nil
}

func GetMyLeaveRequests(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	return repository.GetLeaveRequestsByEmployeeID(employeeID, startDate, endDate)
}

func GetAllCompanyLeaveRequests(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	return repository.GetCompanyLeaveRequestsPaginated(companyID, status, search, startDate, endDate, page, pageSize)
}

func ReviewLeaveRequest(leaveRequestID, adminID uint, status string) (*models.LeaveRequest, error) {
	leaveRequest, err := repository.GetLeaveRequestByID(leaveRequestID)
	if err != nil {
		return nil, fmt.Errorf("leave request not found")
	}

	if leaveRequest == nil {
		return nil, fmt.Errorf("leave request not found")
	}

	employee, err := repository.GetEmployeeByID(int(leaveRequest.EmployeeID))
	if err != nil || employee == nil {
		return nil, fmt.Errorf("could not find employee for leave request")
	}

	adminCompany, err := repository.GetAdminCompanyByID(int(adminID))
	if err != nil || adminCompany == nil || adminCompany.CompanyID != employee.CompanyID {
		return nil, fmt.Errorf("you are not authorized to review this leave request")
	}

	// Only pending requests can be reviewed (approved, rejected, or cancelled by admin)
	if leaveRequest.Status != "pending" {
		return nil, fmt.Errorf("only pending leave requests can be reviewed")
	}

	leaveRequest.Status = status
	leaveRequest.ReviewedBy = &adminID
	now := time.Now()
	leaveRequest.ReviewedAt = &now

	if err := repository.UpdateLeaveRequest(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to update leave request status: %w", err)
	}

	return leaveRequest, nil
}

func ExportCompanyLeaveRequests(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	return repository.GetCompanyLeaveRequestsFiltered(companyID, status, search, startDate, endDate)
}

// CancelLeaveRequest allows an employee to cancel their pending leave request.
func CancelLeaveRequest(leaveRequestID uint, employeeID uint) (*models.LeaveRequest, error) {
	leaveRequest, err := repository.GetLeaveRequestByID(leaveRequestID)
	if err != nil {
		return nil, fmt.Errorf("leave request not found")
	}

	if leaveRequest == nil {
		return nil, fmt.Errorf("leave request not found")
	}

	// Verify that the employee trying to cancel is the owner of the request
	if leaveRequest.EmployeeID != employeeID {
		return nil, fmt.Errorf("you are not authorized to cancel this leave request")
	}

	// Only pending requests can be cancelled
	if leaveRequest.Status != "pending" {
		return nil, fmt.Errorf("only pending leave requests can be cancelled")
	}

	leaveRequest.Status = "cancelled"
	leaveRequest.CancelledByActorType = "employee"
	leaveRequest.CancelledByActorID = &employeeID
	if err := repository.UpdateLeaveRequest(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to cancel leave request: %w", err)
	}

	return leaveRequest, nil
}

// AdminCancelApprovedLeave allows an admin to cancel an already approved leave request.
func AdminCancelApprovedLeave(leaveRequestID uint, adminID uint) (*models.LeaveRequest, error) {
	leaveRequest, err := repository.GetLeaveRequestByID(leaveRequestID)
	if err != nil {
		return nil, fmt.Errorf("leave request not found")
	}

	if leaveRequest == nil {
		return nil, fmt.Errorf("leave request not found")
	}

	// Verify that the leave request is approved
	if leaveRequest.Status != "approved" {
		return nil, fmt.Errorf("only approved leave requests can be cancelled by admin")
	}

	// Verify admin authorization (same company as employee)
	employee, err := repository.GetEmployeeByID(int(leaveRequest.EmployeeID))
	if err != nil || employee == nil {
		return nil, fmt.Errorf("could not find employee for leave request")
	}

	adminCompany, err := repository.GetAdminCompanyByID(int(adminID))
	if err != nil || adminCompany == nil || adminCompany.CompanyID != employee.CompanyID {
		return nil, fmt.Errorf("you are not authorized to cancel this leave request")
	}

	leaveRequest.Status = "cancelled"
	leaveRequest.ReviewedBy = &adminID // Mark as reviewed by admin who cancelled it
	now := time.Now()
	leaveRequest.ReviewedAt = &now
	leaveRequest.CancelledByActorType = "admin"
	leaveRequest.CancelledByActorID = &adminID

	if err := repository.UpdateLeaveRequest(leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to cancel approved leave request: %w", err)
	}

	return leaveRequest, nil
}

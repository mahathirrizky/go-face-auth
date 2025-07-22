package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"time"
)

func ApplyLeave(employeeID uint, leaveType, startDateStr, endDateStr, reason string) (*models.LeaveRequest, error) {
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

	_, err = repository.GetEmployeeByID(int(employeeID))
	if err != nil {
		return nil, fmt.Errorf("employee not found")
	}

	leaveRequest := &models.LeaveRequest{
		EmployeeID: employeeID,
		Type:       leaveType,
		StartDate:  startDate,
		EndDate:    endDate,
		Reason:     reason,
		Status:     "pending",
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

	employee, err := repository.GetEmployeeByID(int(leaveRequest.EmployeeID))
	if err != nil || employee == nil {
		return nil, fmt.Errorf("could not find employee for leave request")
	}

	adminCompany, err := repository.GetAdminCompanyByID(int(adminID))
	if err != nil || adminCompany == nil || adminCompany.CompanyID != employee.CompanyID {
		return nil, fmt.Errorf("you are not authorized to review this leave request")
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

package repository

import (
	"go-face-auth/models"
	"time"
)

// LeaveRequestRepository defines the contract for leave_request-related database operations.
type LeaveRequestRepository interface {
	CreateLeaveRequest(leaveRequest *models.LeaveRequest) error
	GetLeaveRequestByID(id uint) (*models.LeaveRequest, error)
	GetAllLeaveRequests() ([]models.LeaveRequest, error)
	GetLeaveRequestsByEmployeeID(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error)
	GetCompanyLeaveRequestsFiltered(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error)
	UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error
	GetRecentLeaveRequestsByCompanyID(companyID int, limit int) ([]models.LeaveRequest, error)
	IsEmployeeOnApprovedLeave(employeeID int, date time.Time) (*models.LeaveRequest, error)
	IsEmployeeOnApprovedLeaveDateRange(employeeID int, startDate, endDate *time.Time) (bool, error)
	GetPendingLeaveRequestsByEmployeeID(employeeID int) ([]models.LeaveRequest, error)
	GetCompanyLeaveRequestsPaginated(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error)
	GetOnLeaveEmployeesCountToday(companyID int) (int64, error)
}

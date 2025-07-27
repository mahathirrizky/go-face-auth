package repository

import (
	"go-face-auth/models"
	"time"
)

// AttendanceRepository defines the contract for attendance-related database operations.
type AttendanceRepository interface {
	CreateAttendance(attendance *models.AttendancesTable) error
	UpdateAttendance(attendance *models.AttendancesTable) error
	GetLatestAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error)
	GetLatestAttendanceForDate(employeeID int, date time.Time) (*models.AttendancesTable, error)
	GetLatestOvertimeAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error)
	GetPresentEmployeesCountToday(companyID int) (int64, error)
	GetAbsentEmployeesCountToday(companyID int) (int64, error)
	GetAttendancesByCompanyID(companyID int) ([]models.AttendancesTable, error)
	GetRecentAttendancesByCompanyID(companyID int, limit int) ([]models.AttendancesTable, error)
	GetRecentOvertimeAttendancesByCompanyID(companyID int, limit int) ([]models.AttendancesTable, error)
	GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error)
	GetCompanyAttendancesFiltered(companyID int, startDate, endDate *time.Time, attendanceType string) ([]models.AttendancesTable, error)
	HasAttendanceForDate(employeeID int, date time.Time) (bool, error)
	HasAttendanceForDateRange(employeeID int, startDate, endDate *time.Time) (bool, error)
	GetCompanyOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error)
	GetTodayAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error)
	GetRecentAttendancesByEmployeeID(employeeID int, limit int) ([]models.AttendancesTable, error)
	GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error)
	GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error)
	GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	GetUnaccountedEmployeesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.EmployeesTable, error)
	GetOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.AttendancesTable, error)
}

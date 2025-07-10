package models

import "time"

// Attendance represents an attendance record for an employee.
type AttendancesTable struct {
	ID           int            `json:"id"`
	EmployeeID   int            `json:"employee_id"`
	Employee     EmployeesTable `gorm:"foreignKey:EmployeeID" json:"employee"`
	CheckInTime  time.Time      `json:"check_in_time"`
	CheckOutTime *time.Time     `json:"check_out_time"` // Use pointer for nullable DATETIME
	OvertimeMinutes int         `json:"overtime_minutes"` // New field for overtime in minutes
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

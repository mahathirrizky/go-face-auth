package models

import "time"

// Attendance represents an attendance record for an employee.
type Attendance struct {
	ID          int       `json:"id"`
	EmployeeID  int       `json:"employee_id"`
	CheckInTime time.Time `json:"check_in_time"`	
	CheckOutTime *time.Time `json:"check_out_time"` // Use pointer for nullable DATETIME
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

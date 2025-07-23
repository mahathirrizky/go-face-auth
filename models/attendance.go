package models

import "time"

// Attendance represents an attendance record for an employee.
type AttendancesTable struct {
	ID                int             `json:"id"`
	EmployeeID        int             `json:"employee_id"`
	Employee          EmployeesTable  `gorm:"foreignKey:EmployeeID" json:"employee"`
	CheckInTime       time.Time       `json:"check_in_time"`
	CheckOutTime      *time.Time      `json:"check_out_time"` // Use pointer for nullable DATETIME
	OvertimeMinutes   int             `json:"overtime_minutes"`
	Status            string          `json:"status"`
	IsCorrection      bool            `json:"is_correction"`
	Notes             string          `json:"notes"`
	CorrectedByAdminID *uint           `json:"corrected_by_admin_id"` // Nullable admin ID
	CorrectedByAdmin  AdminCompaniesTable `gorm:"foreignKey:CorrectedByAdminID" json:"corrected_by_admin"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

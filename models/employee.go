package models

import "time"

// Employee represents an employee in the system.
type EmployeesTable struct {
	ID             int       `json:"id"`
	CompanyID      int       `json:"company_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Role           string    `json:"role"` // New field for role
	Password       string    `json:"-"` // Exclude from JSON output
	EmployeeIDNumber string    `json:"employee_id_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
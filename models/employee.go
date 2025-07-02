package models

import "time"

// Employee represents an employee in the system.
type Employee struct {
	ID             int       `json:"id"`
	CompanyID      int       `json:"company_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	EmployeeIDNumber string    `json:"employee_id_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
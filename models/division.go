package models

import "time"

// Division represents a department or a section within a company.
type DivisionTable struct {
	ID          uint      `gorm:"primaryKey"`
	CompanyID   uint      `gorm:"not null;index"`
	Name        string    `gorm:"not null;size:255"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Associations
	Company   CompaniesTable           `gorm:"foreignKey:CompanyID"`
	Employees []EmployeesTable         `gorm:"foreignKey:DivisionID"`
	Locations []AttendanceLocation `gorm:"many2many:division_locations;"`
	Shifts    []ShiftsTable            `gorm:"many2many:division_shifts;"`
}

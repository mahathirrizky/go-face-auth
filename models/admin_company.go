package models

import "time"

// AdminCompany represents the relationship between a Company and its Admin Employee.
type AdminCompaniesTable struct {
	ID        int      `json:"id" gorm:"primaryKey"`
	CompanyID int      `json:"company_id" gorm:"uniqueIndex"` // Company that this admin manages
	Company   CompaniesTable `gorm:"foreignKey:CompanyID" json:"-"`
	Email     string    `json:"email"`   // Admin Employee's email
	Password  string    `json:"-"`                      // Admin Employee's password
	Role      string    `json:"role"` // New field for role
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
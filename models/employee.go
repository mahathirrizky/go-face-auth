package models

import "time"

// Employee represents an employee in the system.
type EmployeesTable struct {
	ID                 int        `json:"id"`
	CompanyID          int        `json:"company_id"`
	EmployeeIDNumber   string     `json:"employee_id_number"`
	Name               string     `json:"name"`
	Email              string     `gorm:"unique" json:"email"`
	Password           string     `json:"-"`
	Position           string     `json:"position"`
	Role               string     `json:"role"`
	ShiftID            *int       `json:"shift_id"` // Pointer to allow null
	Shift              ShiftsTable `gorm:"foreignKey:ShiftID" json:"shift"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	FaceImages         []FaceImagesTable `gorm:"foreignKey:EmployeeID" json:"face_images"`
}
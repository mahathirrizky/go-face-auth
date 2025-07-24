package models

import (
	"time"

	"gorm.io/gorm"
)

type LeaveRequest struct {
	gorm.Model
	EmployeeID  uint      `json:"employee_id" gorm:"not null"`
	Employee    EmployeesTable  `json:"employee" gorm:"foreignKey:EmployeeID"` // BelongsTo Employee
	Type        string    `json:"type" gorm:"type:varchar(50);not null"` // e.g., "cuti", "sakit"
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	Reason      string    `json:"reason" gorm:"type:text;not null"`
	Status      string    `json:"status" gorm:"type:varchar(50);default:'pending'"` // e.g., "pending", "approved", "rejected"
	ReviewedBy  *uint     `json:"reviewed_by"` // Admin ID who reviewed it
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CancelledByActorType string `json:"cancelled_by_actor_type,omitempty"` // e.g., "employee", "admin"
	CancelledByActorID *uint `json:"cancelled_by_actor_id,omitempty"` // ID of the employee or admin who cancelled it
	SickNotePath string `json:"sick_note_path,omitempty"` // Path to the uploaded sick note file
}

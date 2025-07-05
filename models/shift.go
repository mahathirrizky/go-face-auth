package models

import "time"

// ShiftsTable represents a work shift for a company.
type ShiftsTable struct {
	ID                 int       `json:"id"`
	CompanyID          int       `json:"company_id"`
	Name               string    `json:"name"`
	StartTime          string    `json:"start_time"` // Stored as "HH:MM:SS"
	EndTime            string    `json:"end_time"`   // Stored as "HH:MM:SS"
	GracePeriodMinutes int       `json:"grace_period_minutes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

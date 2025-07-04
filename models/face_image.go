package models

import "time"

// FaceImage represents a face image associated with an employee.
type FaceImagesTable struct {
	ID        int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}
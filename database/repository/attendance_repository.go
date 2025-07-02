package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateAttendance inserts a new attendance record.
func CreateAttendance(attendance *models.Attendance) error {
	result := database.DB.Create(attendance)
	if result.Error != nil {
		log.Printf("Error creating attendance: %v", result.Error)
		return result.Error
	}
	log.Printf("Attendance record created with ID: %d", attendance.ID)
	return nil
}

// UpdateAttendance updates an existing attendance record.
func UpdateAttendance(attendance *models.Attendance) error {
	result := database.DB.Save(attendance)
	if result.Error != nil {
		log.Printf("Error updating attendance record with ID %d: %v", attendance.ID, result.Error)
		return result.Error
	}
	log.Printf("Attendance record with ID %d updated.", attendance.ID)
	return nil
}

// GetLatestAttendanceByEmployeeID retrieves the latest attendance record for an employee.
// It typically looks for an open check-in (check_out_time IS NULL).
func GetLatestAttendanceByEmployeeID(employeeID int) (*models.Attendance, error) {
	var attendance models.Attendance
	result := database.DB.Where("employee_id = ?", employeeID).Order("check_in_time DESC").Limit(1).First(&attendance)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No attendance record found
		}
		log.Printf("Error getting latest attendance for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}

	return &attendance, nil
}

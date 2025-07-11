package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CreateAttendance inserts a new attendance record.
func CreateAttendance(attendance *models.AttendancesTable) error {
	result := database.DB.Create(attendance)
	if result.Error != nil {
		log.Printf("Error creating attendance: %v", result.Error)
		return result.Error
	}
	log.Printf("Attendance record created with ID: %d", attendance.ID)
	return nil
}

// UpdateAttendance updates an existing attendance record.
func UpdateAttendance(attendance *models.AttendancesTable) error {
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
func GetLatestAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	var attendance models.AttendancesTable
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

// GetLatestOvertimeAttendanceByEmployeeID retrieves the latest overtime attendance record for an employee.
func GetLatestOvertimeAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	var attendance models.AttendancesTable
	result := database.DB.Where("employee_id = ? AND (status = ? OR status = ?)", employeeID, "overtime_in", "overtime_out").Order("check_in_time DESC").Limit(1).First(&attendance)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No overtime attendance record found
		}
		log.Printf("Error getting latest overtime attendance for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}

	return &attendance, nil
}

// GetPresentEmployeesCountToday retrieves the count of employees marked as 'present' for a given company today.
func GetPresentEmployeesCountToday(companyID int) (int64, error) {
	var count int64
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	result := database.DB.Model(&models.AttendancesTable{}).Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ? AND attendances_tables.status = ? AND attendances_tables.check_in_time >= ? AND attendances_tables.check_in_time < ?", companyID, "present", startOfDay, endOfDay).Count(&count)
	if result.Error != nil {
		log.Printf("Error getting present employees count today for company %d: %v", companyID, result.Error)
		return 0, result.Error
	}
	return count, nil
}

// GetAbsentEmployeesCountToday retrieves the count of employees marked as 'absent' for a given company today.
func GetAbsentEmployeesCountToday(companyID int) (int64, error) {
	var count int64
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	result := database.DB.Model(&models.AttendancesTable{}).Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ? AND attendances_tables.status = ? AND attendances_tables.check_in_time >= ? AND attendances_tables.check_in_time < ?", companyID, "absent", startOfDay, endOfDay).Count(&count)
	if result.Error != nil {
		log.Printf("Error getting absent employees count today for company %d: %v", companyID, result.Error)
		return 0, result.Error
	}
	return count, nil
}



// GetAttendancesByCompanyID retrieves all attendance records for a given company ID.
func GetAttendancesByCompanyID(companyID int) ([]models.AttendancesTable, error) {
	var attendances []models.AttendancesTable
	result := database.DB.Preload("Employee").Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ?", companyID).Order("attendances_tables.check_in_time desc").Find(&attendances)
	if result.Error != nil {
		log.Printf("Error querying attendances for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return attendances, nil
}

// GetRecentAttendancesByCompanyID retrieves recent attendance records for a given company ID.
func GetRecentAttendancesByCompanyID(companyID int, limit int) ([]models.AttendancesTable, error) {
	var attendances []models.AttendancesTable
	result := database.DB.Preload("Employee").Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ?", companyID).Order("check_in_time DESC").Limit(limit).Find(&attendances)
	if result.Error != nil {
		log.Printf("Error getting recent attendances for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return attendances, nil
}

// GetEmployeeAttendances retrieves attendance records for a specific employee, optionally filtered by date range.
func GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	var attendances []models.AttendancesTable
	query := database.DB.Preload("Employee").Where("employee_id = ?", employeeID)

	if startDate != nil {
		query = query.Where("check_in_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("check_in_time <= ?", *endDate)
	}

	result := query.Order("check_in_time DESC").Find(&attendances)
	if result.Error != nil {
		log.Printf("Error getting attendance for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return attendances, nil
}

// GetCompanyAttendancesFiltered retrieves all attendance records for a given company ID, optionally filtered by date range.
func GetCompanyAttendancesFiltered(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	var attendances []models.AttendancesTable
	query := database.DB.Preload("Employee").Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ?", companyID)

	if startDate != nil {
		query = query.Where("attendances_tables.check_in_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("attendances_tables.check_in_time <= ?", *endDate)
	}

	result := query.Order("attendances_tables.check_in_time desc").Find(&attendances)
	if result.Error != nil {
		log.Printf("Error querying attendances for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return attendances, nil
}

// HasAttendanceForDate checks if an employee has any attendance record for a specific date.
func HasAttendanceForDate(employeeID int, date time.Time) (bool, error) {
	var count int64
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second) // End of day

	result := database.DB.Model(&models.AttendancesTable{}).Where("employee_id = ? AND check_in_time >= ? AND check_in_time <= ?", employeeID, startOfDay, endOfDay).Count(&count)
	if result.Error != nil {
		log.Printf("Error checking attendance for employee %d on date %s: %v", employeeID, date.Format("2006-01-02"), result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

// HasAttendanceForDateRange checks if an employee has any attendance record within a specific date range.
func HasAttendanceForDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	var count int64
	query := database.DB.Model(&models.AttendancesTable{}).Where("employee_id = ?", employeeID)

	if startDate != nil {
		query = query.Where("check_in_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("check_in_time <= ?", *endDate)
	}

	result := query.Count(&count)
	if result.Error != nil {
		log.Printf("Error checking attendance for employee %d in date range: %v", employeeID, result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

// GetCompanyOvertimeAttendancesFiltered retrieves all overtime attendance records for a given company ID, optionally filtered by date range.
func GetCompanyOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	var attendances []models.AttendancesTable
	query := database.DB.Preload("Employee").Joins("join employees_tables on employees_tables.id = attendances_tables.employee_id").Where("employees_tables.company_id = ? AND (attendances_tables.status = ? OR attendances_tables.status = ?)", companyID, "overtime_in", "overtime_out")

	if startDate != nil {
		query = query.Where("attendances_tables.check_in_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("attendances_tables.check_in_time <= ?", *endDate)
	}

	result := query.Order("attendances_tables.check_in_time desc").Find(&attendances)
	if result.Error != nil {
		log.Printf("Error querying overtime attendances for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return attendances, nil
}

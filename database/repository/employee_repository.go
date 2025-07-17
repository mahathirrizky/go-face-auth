package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CreateEmployee inserts a new employee into the database.
func CreateEmployee(employee *models.EmployeesTable) error {
	employee.IsPasswordSet = false // Set to false by default for new employees
	result := database.DB.Create(employee)
	if result.Error != nil {
		log.Printf("Error creating employee: %v", result.Error)
		return result.Error
	}
	log.Printf("Employee created with ID: %d", employee.ID)
	return nil
}

// GetEmployeeByID retrieves an employee by their ID.
func GetEmployeeByID(id int) (*models.EmployeesTable, error) {
	var employee models.EmployeesTable
	result := database.DB.Preload("Shift").First(&employee, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Employee not found
		}
		log.Printf("Error getting employee with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &employee, nil
}

// GetEmployeesByCompanyID retrieves all employees for a given company ID.
func GetEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error) {
	var employees []models.EmployeesTable
	result := database.DB.Where("company_id = ?", companyID).Find(&employees)
	if result.Error != nil {
		log.Printf("Error querying employees for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return employees, nil
}

// GetEmployeeByEmail retrieves an employee by their email address.
func GetEmployeeByEmail(email string) (*models.EmployeesTable, error) {
	var employee models.EmployeesTable
	result := database.DB.Where("email = ?", email).First(&employee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Employee not found
		}
		log.Printf("Error getting employee with email %s: %v", email, result.Error)
		return nil, result.Error
	}
	return &employee, nil
}

// UpdateEmployee updates an existing employee record in the database.
func UpdateEmployee(employee *models.EmployeesTable) error {
	result := database.DB.Save(employee)
	if result.Error != nil {
		log.Printf("Error updating employee: %v", result.Error)
		return result.Error
	}
	log.Printf("Employee updated with ID: %d", employee.ID)
	return nil
}

// DeleteEmployee removes an employee from the database by their ID.
func DeleteEmployee(id int) error {
	log.Printf("Attempting to delete employee with ID: %d", id)
	result := database.DB.Delete(&models.EmployeesTable{}, id)
	if result.Error != nil {
		log.Printf("Error deleting employee with ID %d: %v", id, result.Error)
		return result.Error
	}
	log.Printf("Delete operation for employee ID %d. Rows affected: %d", id, result.RowsAffected)
	if result.RowsAffected == 0 {
		log.Printf("No employee found with ID %d to delete or already deleted", id)
		return gorm.ErrRecordNotFound // Or a custom error
	}
	log.Printf("Employee with ID %d deleted successfully", id)
	return nil
}

// SearchEmployees finds employees by name within a specific company.
func SearchEmployees(companyID int, name string) ([]models.EmployeesTable, error) {
	var employees []models.EmployeesTable
	// Use ILIKE for case-insensitive search, works on PostgreSQL. Use LIKE for others.
	result := database.DB.Where("company_id = ? AND LOWER(name) LIKE LOWER(?)", companyID, "%"+name+"%").Find(&employees)
	if result.Error != nil {
		log.Printf("Error searching for employees with name %s in company %d: %v", name, companyID, result.Error)
		return nil, result.Error
	}
	return employees, nil
}

// GetTotalEmployeesByCompanyID retrieves the total number of employees for a given company ID.
func GetTotalEmployeesByCompanyID(companyID int) (int64, error) {
	var count int64
	result := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", companyID).Count(&count)
	if result.Error != nil {
		log.Printf("Error counting employees for company %d: %v", companyID, result.Error)
		return 0, result.Error
	}
	return count, nil
}

// GetEmployeesWithFaceImages retrieves all employees for a company, preloading their face images.
func GetEmployeesWithFaceImages(companyID int) ([]models.EmployeesTable, error) {
	var employees []models.EmployeesTable
	result := database.DB.Preload("FaceImages").Where("company_id = ?", companyID).Find(&employees)
	if result.Error != nil {
		log.Printf("Error getting employees with face images for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return employees, nil
}

// GetOnLeaveEmployeesCountToday retrieves the count of employees who are on leave today for a given company.
func GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD for date comparison

	result := database.DB.Model(&models.LeaveRequest{}).
		Joins("JOIN employees_tables ON leave_requests.employee_id = employees_tables.id").
		Where("employees_tables.company_id = ? AND leave_requests.status = ? AND ? BETWEEN leave_requests.start_date AND leave_requests.end_date",
			companyID, "approved", today).
		Count(&count)

	if result.Error != nil {
		log.Printf("Error counting on-leave employees today for company %d: %v", companyID, result.Error)
		return 0, result.Error
	}
	return count, nil
}

// SetEmployeePasswordSet updates the IsPasswordSet field for an employee.
func SetEmployeePasswordSet(employeeID uint, isSet bool) error {
	result := database.DB.Model(&models.EmployeesTable{}).Where("id = ?", employeeID).Update("is_password_set", isSet)
	if result.Error != nil {
		log.Printf("Error setting IsPasswordSet for employee %d: %v", employeeID, result.Error)
		return result.Error
	}
	return nil
}

// GetPendingEmployees retrieves all employees who have not set their password yet.
func GetPendingEmployees(companyID int) ([]models.EmployeesTable, error) {
	var employees []models.EmployeesTable
	result := database.DB.Where("company_id = ? AND is_password_set = ?", companyID, false).Find(&employees)
	if result.Error != nil {
		log.Printf("Error getting pending employees for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return employees, nil
}

// GetEmployeeByEmailOrIDNumber retrieves an employee by their email address or employee ID number.
func GetEmployeeByEmailOrIDNumber(email, employeeIDNumber string) (*models.EmployeesTable, error) {
	var employee models.EmployeesTable
	result := database.DB.Where("email = ? OR employee_id_number = ?", email, employeeIDNumber).First(&employee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound // Return a specific error for not found
		}
		log.Printf("Error getting employee by email %s or ID number %s: %v", email, employeeIDNumber, result.Error)
		return nil, result.Error
	}
	return &employee, nil
}

// GetEmployeesByCompanyIDPaginated retrieves paginated employees for a given company ID.
func GetEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	var employees []models.EmployeesTable
	var totalRecords int64

	query := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ?", companyID)

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR employee_id_number ILIKE ? OR position ILIKE ?", searchQuery, searchQuery, searchQuery, searchQuery)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Preload("Shift").Offset(offset).Limit(pageSize).Find(&employees)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return employees, totalRecords, nil
}

// GetPendingEmployeesByCompanyIDPaginated retrieves paginated pending employees for a given company ID.
func GetPendingEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	var employees []models.EmployeesTable
	var totalRecords int64

	query := database.DB.Model(&models.EmployeesTable{}).Where("company_id = ? AND is_password_set = ?", companyID, false)

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Find(&employees)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return employees, totalRecords, nil
}
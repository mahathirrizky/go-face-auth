package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateEmployee inserts a new employee into the database.
func CreateEmployee(employee *models.EmployeesTable) error {
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
	result := database.DB.First(&employee, id)
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

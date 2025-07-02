package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateEmployee inserts a new employee into the database.
func CreateEmployee(employee *models.Employee) error {
	result := database.DB.Create(employee)
	if result.Error != nil {
		log.Printf("Error creating employee: %v", result.Error)
		return result.Error
	}
	log.Printf("Employee created with ID: %d", employee.ID)
	return nil
}

// GetEmployeeByID retrieves an employee by their ID.
func GetEmployeeByID(id int) (*models.Employee, error) {
	var employee models.Employee
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
func GetEmployeesByCompanyID(companyID int) ([]models.Employee, error) {
	var employees []models.Employee
	result := database.DB.Where("company_id = ?", companyID).Find(&employees)
	if result.Error != nil {
		log.Printf("Error querying employees for company %d: %v", companyID, result.Error)
		return nil, result.Error
	}
	return employees, nil
}

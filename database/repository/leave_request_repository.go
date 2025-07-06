package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateLeaveRequest inserts a new leave request into the database.
func CreateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := database.DB.Create(leaveRequest)
	if result.Error != nil {
		log.Printf("Error creating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request created with ID: %d", leaveRequest.ID)
	return nil
}

// GetLeaveRequestByID retrieves a leave request by its ID.
func GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	var leaveRequest models.LeaveRequest
	result := database.DB.Preload("Employee").First(&leaveRequest, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Leave request not found
		}
		log.Printf("Error getting leave request with ID %d: %v", id, result.Error)
		return nil, result.Error
	}
	return &leaveRequest, nil
}

// GetAllLeaveRequests retrieves all leave requests.
func GetAllLeaveRequests() ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := database.DB.Preload("Employee").Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting all leave requests: %v", result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// GetLeaveRequestsByEmployeeID retrieves all leave requests for a given employee ID.
func GetLeaveRequestsByEmployeeID(employeeID uint) ([]models.LeaveRequest, error) {
	var leaveRequests []models.LeaveRequest
	result := database.DB.Preload("Employee").Where("employee_id = ?", employeeID).Find(&leaveRequests)
	if result.Error != nil {
		log.Printf("Error getting leave requests for employee %d: %v", employeeID, result.Error)
		return nil, result.Error
	}
	return leaveRequests, nil
}

// UpdateLeaveRequest updates an existing leave request record in the database.
func UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	result := database.DB.Save(leaveRequest)
	if result.Error != nil {
		log.Printf("Error updating leave request: %v", result.Error)
		return result.Error
	}
	log.Printf("Leave request updated with ID: %d", leaveRequest.ID)
	return nil
}

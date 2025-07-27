package repository

import "go-face-auth/models"

// BroadcastRepository defines the contract for broadcast-related database operations.
type BroadcastRepository interface {
	CreateBroadcast(message *models.BroadcastMessage) error
	GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error)
	MarkBroadcastAsRead(employeeID, messageID uint) error
}

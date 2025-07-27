package repository

import (
	"go-face-auth/models"
	"time"

	"gorm.io/gorm"
)

type broadcastRepository struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) BroadcastRepository {
	return &broadcastRepository{db: db}
}

// Create saves a new broadcast message to the database.
func (r *broadcastRepository) CreateBroadcast(message *models.BroadcastMessage) error {
	return r.db.Create(message).Error
}

// GetForEmployee retrieves all active broadcast messages for a company,
// and marks them as read if the employee has read them.
func (r *broadcastRepository) GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error) {
	var messages []models.BroadcastMessage

	// This complex query does the following:
	// 1. Selects all columns from `broadcast_messages`.
	// 2. Adds a boolean column `is_read` which is true if a corresponding entry exists in `employee_broadcast_reads`.
	// 3. Filters messages for the correct company_id.
	// 4. Filters out messages that are expired (expire_date is in the past).
	// 5. Orders by creation date, newest first.
	err := r.db.Table("broadcast_messages").
		Select("broadcast_messages.*, CASE WHEN ebr.employee_id IS NOT NULL THEN TRUE ELSE FALSE END as is_read").
		Joins("LEFT JOIN employee_broadcast_reads ebr ON ebr.broadcast_message_id = broadcast_messages.id AND ebr.employee_id = ?", employeeID).
		Where("broadcast_messages.company_id = ?", companyID).
		Where("broadcast_messages.expire_date IS NULL OR broadcast_messages.expire_date > ?", time.Now()).
		Order("broadcast_messages.created_at DESC").
		Find(&messages).Error

	return messages, err
}

// MarkAsRead creates a record indicating an employee has read a message.
func (r *broadcastRepository) MarkBroadcastAsRead(employeeID, messageID uint) error {
	read := models.EmployeeBroadcastRead{
		EmployeeID:        employeeID,
		BroadcastMessageID: messageID,
		ReadAt:            time.Now(),
	}
	// Using FirstOrCreate to prevent duplicate entries if the request is sent multiple times.
	result := r.db.FirstOrCreate(&read)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
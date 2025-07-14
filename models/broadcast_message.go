package models

import "time"

// BroadcastMessage represents a message sent by an admin to a company.
type BroadcastMessage struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CompanyID   uint       `gorm:"not null" json:"company_id"`
	Message     string     `gorm:"type:text;not null" json:"message"`
	ExpireDate  *time.Time `gorm:"null" json:"expire_date"`
	CreatedAt   time.Time  `json:"created_at"`
	IsRead      bool       `gorm:"-" json:"is_read"` // This field is populated by a custom query, not a DB column.
}

// EmployeeBroadcastRead tracks which employee has read which broadcast message.
type EmployeeBroadcastRead struct {
	EmployeeID        uint      `gorm:"primaryKey" json:"employee_id"`
	BroadcastMessageID uint      `gorm:"primaryKey" json:"broadcast_message_id"`
	ReadAt            time.Time `json:"read_at"`
}

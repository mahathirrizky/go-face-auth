package services

import (
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"time"
)

func CreateBroadcastMessage(companyID uint, message string, expireDate string) (*models.BroadcastMessage, error) {
	var expireTime *time.Time
	if expireDate != "" {
		parsedTime, err := time.Parse("2006-01-02", expireDate)
		if err != nil {
			return nil, err
		}
		expireTime = &parsedTime
	}

	newMessage := &models.BroadcastMessage{
		CompanyID:  companyID,
		Message:    message,
		ExpireDate: expireTime,
	}

	if err := repository.CreateBroadcast(newMessage); err != nil {
		return nil, err
	}

	return newMessage, nil
}

func GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error) {
	return repository.GetBroadcastsForEmployee(companyID, employeeID)
}

func MarkBroadcastAsRead(employeeID, messageID uint) error {
	return repository.MarkBroadcastAsRead(employeeID, messageID)
}

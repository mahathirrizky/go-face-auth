package services

import (
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"time"
)

// BroadcastService defines the interface for broadcast related business logic.
type BroadcastService interface {
	CreateBroadcastMessage(companyID uint, message string, expireDate string) (*models.BroadcastMessage, error)
	GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error)
	MarkBroadcastAsRead(employeeID, messageID uint) error
}

// broadcastService is the concrete implementation of BroadcastService.
type broadcastService struct {
	broadcastRepo repository.BroadcastRepository
}

// NewBroadcastService creates a new instance of BroadcastService.
func NewBroadcastService(broadcastRepo repository.BroadcastRepository) BroadcastService {
	return &broadcastService{
		broadcastRepo: broadcastRepo,
	}
}

func (s *broadcastService) CreateBroadcastMessage(companyID uint, message string, expireDate string) (*models.BroadcastMessage, error) {
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

	if err := s.broadcastRepo.CreateBroadcast(newMessage); err != nil {
		return nil, err
	}

	return newMessage, nil
}

func (s *broadcastService) GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error) {
	return s.broadcastRepo.GetBroadcastsForEmployee(companyID, employeeID)
}

func (s *broadcastService) MarkBroadcastAsRead(employeeID, messageID uint) error {
	return s.broadcastRepo.MarkBroadcastAsRead(employeeID, messageID)
}
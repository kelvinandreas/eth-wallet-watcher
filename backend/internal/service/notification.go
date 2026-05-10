package service

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) GetByUserID(userID uuid.UUID) ([]app.Notification, error) {
	return s.notifRepo.GetByUserID(userID)
}

func (s *NotificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	return s.notifRepo.MarkAsRead(userID, notificationID)
}

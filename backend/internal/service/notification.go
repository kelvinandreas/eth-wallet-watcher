package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	notifRepo   *repository.NotificationRepository
	redisClient *redis.Client
}

func NewNotificationService(notifRepo *repository.NotificationRepository, redisClient *redis.Client) *NotificationService {
	return &NotificationService{notifRepo: notifRepo, redisClient: redisClient}
}

func (s *NotificationService) GetByUserID(userID uuid.UUID, page, limit int) ([]app.Notification, int64, error) {
	key := helper.CacheKey(constant.CacheKey.Notifications, userID)

	all, err := helper.GetOrSet(context.Background(), s.redisClient, key, helper.DefaultCacheTTL, func() ([]app.Notification, error) {
		return s.notifRepo.GetByUserID(userID)
	})
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(all))
	start, end := paginate(len(all), page, limit)
	return all[start:end], total, nil
}

func (s *NotificationService) MarkAsRead(userID, notificationID uuid.UUID) error {
	if err := s.notifRepo.MarkAsRead(userID, notificationID); err != nil {
		return err
	}
	key := helper.CacheKey(constant.CacheKey.Notifications, userID)
	helper.DeleteCache(context.Background(), s.redisClient, key)
	return nil
}

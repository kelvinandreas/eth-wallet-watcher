package repository

import (
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *app.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) GetByUserID(userID uuid.UUID) ([]app.Notification, error) {
	var notifications []app.Notification
	result := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications)
	return notifications, result.Error
}

func (r *NotificationRepository) MarkAsRead(userID, notificationID uuid.UUID) error {
	result := r.db.Model(&app.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/helper"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/response"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/service"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	notifService *service.NotificationService
}

func NewNotificationHandler(notifService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	notifications, err := h.notifService.GetByUserID(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	result := make([]response.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, response.NewNotificationResponse(n))
	}

	return response.Success(c, fiber.StatusOK, constant.MsgNotificationsRetrieved, result)
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notifID, err := uuid.Parse(strings.TrimSpace(c.Params("notifID")))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, constant.ErrInvalidNotificationID)
	}

	userID, err := helper.GetUserIDFromLocals(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	if err := h.notifService.MarkAsRead(userID, notifID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, constant.ErrNotificationNotFound)
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, constant.MsgNotificationRead, nil)
}

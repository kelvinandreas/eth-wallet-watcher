package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
)

type NotificationResponse struct {
	ID       uuid.UUID `json:"id"`
	WalletID uuid.UUID `json:"wallet_id"`
	TxHash   string    `json:"tx_hash"`
	Message  string    `json:"message"`
	IsRead   bool      `json:"is_read"`
	DateIn   time.Time `json:"date_in"`
}

func NewNotificationResponse(n app.Notification) NotificationResponse {
	return NotificationResponse{
		ID:       n.ID,
		WalletID: n.WalletID,
		TxHash:   n.TxHash,
		Message:  n.Message,
		IsRead:   n.IsRead,
		DateIn:   n.DateIn,
	}
}

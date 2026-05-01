package app

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	BaseModel
	WalletID    uuid.UUID       `gorm:"type:uuid;not null;index"`
	Wallet      MonitoredWallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TxHash      string          `gorm:"size:66;not null;uniqueIndex"`
	FromAddress string          `gorm:"size:42;not null;index"`
	ToAddress   string          `gorm:"size:42;not null;index"`
	Value       string          `gorm:"type:numeric(78,0);not null"`
	Timestamp   time.Time       `gorm:"not null;index"`
}

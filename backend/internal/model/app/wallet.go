package app

import (
	"github.com/google/uuid"
)

type MonitoredWallet struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Address   string    `gorm:"size:42;not null;index"`
	Label     string    `gorm:"size:120"`
	LastBlock string    `gorm:"column:last_block"`
}

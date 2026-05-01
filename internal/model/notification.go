package model

import "github.com/google/uuid"

type Notification struct {
	BaseModel
	UserID    uuid.UUID       `gorm:"type:uuid;not null;index"`
	User      User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID  uuid.UUID       `gorm:"type:uuid;not null;index"`
	Wallet    MonitoredWallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TxHash    string          `gorm:"size:66;not null;index"`
	Message   string          `gorm:"type:text;not null"`
	IsRead    bool            `gorm:"not null;default:false"`
	EmailSent bool            `gorm:"not null;default:false"`
}

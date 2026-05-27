package app

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateIn time.Time `gorm:"column:created_at;autoCreateTime"`
	DateUp time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	return nil
}

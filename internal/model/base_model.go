package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID        uuid.UUID             `gorm:"type:uuid;primaryKey"`
	DateIn    time.Time             `gorm:"column:created_at;autoCreateTime"`
	DateUp    time.Time             `gorm:"column:updated_at;autoUpdateTime"`
	UserIn    *uuid.UUID            `gorm:"type:uuid"`
	UserUp    *uuid.UUID            `gorm:"type:uuid"`
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	return nil
}

package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModelUUID struct {
	ID string `gorm:"type:char(36);primaryKey"`
}

func (b *BaseModelUUID) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

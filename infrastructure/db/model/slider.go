package model

import (
	"time"

	"gorm.io/gorm"
)

type Slider struct {
	BaseModelUUID
	Title       string         `gorm:"not null"`
	Description string         `gorm:"not null"`
	File        string         `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

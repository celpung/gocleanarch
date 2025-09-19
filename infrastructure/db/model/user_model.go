package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModelUUID
	Name      string
	Email     string         `gorm:"unique"`
	Password  string         `gorm:"not null"`
	Active    bool           `gorm:"default:0"`
	Role      string         `gorm:"not null;default:1"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

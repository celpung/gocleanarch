package model

import "time"

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey;column:id"`
	Name      string    `gorm:"size:120;not null;column:name"`
	Email     string    `gorm:"size:120;not null;uniqueIndex;column:email"`
	Password  string    `gorm:"size:255;not null;column:password"`
	Role      string    `gorm:"size:50;not null;default:user;column:role"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"password"`
	Role      int32      `json:"role"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

// type UserRepository interface {
// 	Create(user *User) error
// 	Read(ctx context.Context) (user []User, err error)
// }

type UserUsecase interface {
	Read(ctx context.Context) ([]User, error)
}

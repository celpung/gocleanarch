package user_entity

import "time"

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name" binding:"required"`
	Email     string     `gorm:"unique" json:"email" binding:"required"`
	Password  string     `gorm:"not null" json:"password" binding:"required,min=8"`
	Active    bool       `gorm:"default:0" json:"active"`
	Role      uint       `gorm:"not null;default:1" json:"role"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type UserUpdate struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
	Role     uint   `json:"role"`
}

type UserHttpResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
	Role   uint   `json:"role"`
}

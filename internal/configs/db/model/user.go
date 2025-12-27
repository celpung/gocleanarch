package model

type User struct {
	ID       string `gorm:"type:char(36);primaryKey;column:id"`
	Name     string `gorm:"size:120;not null;column:name"`
	Email    string `gorm:"size:120;not null;uniqueIndex;column:email"`
	Password string `gorm:"size:255;not null;column:password"`
	IsActive bool   `gorm:"not null;default:true;column:is_active"`
}

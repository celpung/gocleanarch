package slider_entity

import "time"

type Slider struct {
	ID          uint       `gorm:"primaryKey" json:"id" form:"id"`
	Title       string     `json:"title" binding:"required" form:"title"`
	Description string     `json:"description" binding:"required" form:"description"`
	File        string     `json:"file" binding:"required" form:"file"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

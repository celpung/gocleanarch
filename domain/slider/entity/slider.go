package entity

import "time"

type Slider struct {
	ID          string
	Title       string
	Description string
	File        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

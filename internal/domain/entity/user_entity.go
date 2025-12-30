package entity

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUser struct {
	Name     *string
	Email    *string
	Role     *string
	Password *string
}

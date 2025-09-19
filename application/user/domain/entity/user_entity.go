package entity

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Active    bool
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type UpdateUserPayload struct {
	ID       string
	Name     *string
	Email    *string
	Password *string
	Active   *bool
	Role     *string
}

package dto

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	ID       string  `json:"id" validate:"required"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Name     *string `json:"name" validate:"omitempty,min=2"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	Role     *string `json:"role" validate:"omitempty,min=2"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

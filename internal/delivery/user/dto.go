package user

import "time"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Role     *string `json:"role" validate:"omitempty,min=2"`
	Password *string `json:"password" validate:"omitempty,min=8"`
}

type PagingMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Data []UserResponse `json:"data"`
	Meta PagingMeta     `json:"meta"`
}

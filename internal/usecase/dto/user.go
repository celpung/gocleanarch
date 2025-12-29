package dto

import "time"

// CreateUserRequest represents user creation input at the usecase boundary.
type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
	Role     string
}

// UpdateUserRequest represents partial updates for a user.
type UpdateUserRequest struct {
	Email    *string
	Name     *string
	Password *string
	Role     *string
}

// UserResponse is the data returned by usecases to delivery.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package dto

import userdto "github.com/celpung/gocleanarch/internal/usecase/user/dto"

type RegisterRequest = userdto.CreateUserRequest

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

type UserResponse = userdto.UserResponse

type ListUsersResponse struct {
	Data []UserResponse `json:"data"`
	Meta PagingMeta     `json:"meta"`
}

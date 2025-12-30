package dto

import (
	"time"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	"github.com/celpung/gocleanarch/pkg/utilities/mapper"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func ToUserResponse(u entity.User) (UserResponse, error) {
	var resp UserResponse
	if err := mapper.CopyTo(&u, &resp); err != nil {
		return UserResponse{}, err
	}
	return resp, nil
}

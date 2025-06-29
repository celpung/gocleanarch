package dto

import "github.com/celpung/gocleanarch/domain/user/entity"

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type UserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
	Role   uint   `json:"role"`
}

func UserResponseDTO(entity *entity.User) *UserResponse {
	return &UserResponse{
		ID:     entity.ID,
		Name:   entity.Name,
		Email:  entity.Email,
		Active: entity.Active,
		Role:   entity.Role,
	}
}

func UserResponseListDTO(entities []*entity.User) []*UserResponse {
	var responses []*UserResponse
	for _, entity := range entities {
		responses = append(responses, UserResponseDTO(entity))
	}
	return responses
}

// DTO untuk create
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	Role     uint   `json:"role" binding:"required" validate:"required"`
}

// DTO untuk update
type UserUpdateRequest struct {
	ID       uint   `json:"id" binding:"required" validate:"required"`
	Name     string `json:"name" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"omitempty,min=8" validate:"omitempty,min=8"`
	Active   bool   `json:"active" binding:"omitempty" validate:"omitempty"`
	Role     uint   `json:"role" binding:"omitempty" validate:"omitempty"`
}

// Konversi Create
func UserCreateRequestDTO(dto *UserCreateRequest) *entity.User {
	return &entity.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
		Active:   true,
	}
}

// Konversi Update
func UserUpdateRequestDTO(dto *UserUpdateRequest) *entity.User {
	return &entity.User{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Active:   dto.Active,
		Role:     dto.Role,
	}
}

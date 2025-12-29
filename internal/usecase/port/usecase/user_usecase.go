package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/usecase/dto"
)

type UserUsecase interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID, password string) error
	GetByID(ctx context.Context, id string) (dto.UserResponse, error)
	List(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error)
	Update(ctx context.Context, id string, updates dto.UpdateUserRequest) (dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

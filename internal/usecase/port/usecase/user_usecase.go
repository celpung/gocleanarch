package usecase

import (
	"context"

	userdto "github.com/celpung/gocleanarch/internal/usecase/user/dto"
)

type UserUsecase interface {
	Create(ctx context.Context, req userdto.CreateUserRequest) (userdto.UserResponse, error)
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID, password string) error
	GetByID(ctx context.Context, id string) (userdto.UserResponse, error)
	List(ctx context.Context, page, limit int) ([]userdto.UserResponse, int64, error)
	Update(ctx context.Context, req userdto.UpdateUserRequest) (userdto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

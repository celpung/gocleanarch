package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/entity"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*entity.User, error)
	GetUser(ctx context.Context, id string) (*entity.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) error
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type CreateUserInput struct {
	Name      string
	Email     string
	Password  string
	CompanyID *string
}

type UpdateUserInput struct {
	Name      *string
	Email     *string
	Password  *string
	CompanyID *string
	IsActive  *bool
}

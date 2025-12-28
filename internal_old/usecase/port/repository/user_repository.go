package repository

import (
	"context"

	"github.com/celpung/gocleanarch/internal_old/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	ListUsers(ctx context.Context, offset, limit int) ([]entity.User, int64, error)
	UpdateUser(ctx context.Context, id string, fields map[string]any) error
	DeleteUser(ctx context.Context, id string) error
}

package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/entity"
)

type UserUsecase interface {
	Register(ctx context.Context, user entity.User) error
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, email, password string) error
	UserLists(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	UpdateUser(ctx context.Context, id string, input *entity.User) error
	Delete(ctx context.Context, id string) error
}

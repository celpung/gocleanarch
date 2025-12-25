package port

import (
	"context"

	"github.com/celpung/gocleanarch/internal/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/internal/modules/user/usecase/dto"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, payload *entity.User) error
	Authenticate(ctx context.Context, email, password string) (*entity.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	UpdateUser(ctx context.Context, id string, payload dto.UpdateUserInput) error
	DeleteUser(ctx context.Context, id string) error
}

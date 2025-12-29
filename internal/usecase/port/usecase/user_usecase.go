package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	"github.com/celpung/gocleanarch/internal/usecase/dto"
)

type UserUsecase interface {
	Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID, password string) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	List(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	Update(ctx context.Context, id string, updates dto.UpdateUserInput) (entity.User, error)
	Delete(ctx context.Context, id string) error
}

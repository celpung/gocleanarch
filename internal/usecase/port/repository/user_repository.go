package repository

import (
	"context"

	"github.com/celpung/gocleanarch/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	Lists(ctx context.Context, offset, limit int) ([]entity.User, int64, error)
	FindByID(ctx context.Context, userID string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, id string, input *entity.UpdateUser) error
	Delete(ctx context.Context, id string) error
}

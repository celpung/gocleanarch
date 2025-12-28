package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/entity"
)

type UserUsecase interface {
	Register(ctx context.Context, user entity.User) error
	Login(ctx context.Context, email, password string) (string, error)

	// Create(ctx context.Context, user entity.User) error
	// Lists(ctx context.Context, offset, limit int) ([]entity.User, int64, error)
	// FindByID(ctx context.Context, userID string) (*entity.User, error)
	// FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// Update(ctx context.Context, id string, fields map[string]any) error
	// Delete(ctx context.Context, id string) error
}

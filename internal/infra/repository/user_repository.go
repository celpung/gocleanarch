package repository

import (
	"context"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	panic("unimplemented")
}

func (r *UserRepository) Lists(ctx context.Context, offset int, limit int) ([]entity.User, int64, error) {
	panic("unimplemented")
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*entity.User, error) {
	panic("unimplemented")
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	panic("unimplemented")
}

func (r *UserRepository) Update(ctx context.Context, id string, input *entity.UpdateUser) error {
	panic("unimplemented")
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{
		db: db,
	}
}

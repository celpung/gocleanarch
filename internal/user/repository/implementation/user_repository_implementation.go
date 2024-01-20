package repositoryimplementation

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
	"github.com/celpung/gocleanarch/internal/user/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Read(ctx context.Context) ([]domain.User, error) {
	var user []domain.User
	err := r.DB.Find(&user).Error
	return user, err
}

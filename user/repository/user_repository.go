package repository

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	// Conn *sql.DB
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

// Read implements the Read method of the domain.UserRepository interface.
func (r *UserRepository) Read(ctx context.Context) ([]domain.User, error) {
	var user []domain.User
	err := r.DB.Find(&user).Error
	return user, err
}

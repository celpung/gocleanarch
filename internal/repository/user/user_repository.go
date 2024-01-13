package repository

import (
	"github.com/celpung/gocleanarch/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create implements UserRepositoryInterface.
func (r *UserRepository) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

// Read implements UserRepositoryInterface.
func (r *UserRepository) Read() ([]entity.User, error) {
	var user []entity.User
	err := r.DB.Find(&user).Error
	return user, err
}

var _ UserInterface = &UserRepository{} // <-- check for the interface implementation

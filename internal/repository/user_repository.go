package repository

import "github.com/celpung/gocleanarch/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	Read() ([]entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}

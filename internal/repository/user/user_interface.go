package repository

import (
	"github.com/celpung/gocleanarch/internal/entity"
)

type UserInterface interface {
	Create(user *entity.User) error
	Read() ([]entity.User, error)
}

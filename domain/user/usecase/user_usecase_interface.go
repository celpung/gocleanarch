package user_usecase

import "github.com/celpung/gocleanarch/entity"

type UserUsecaseInterface interface {
	Create(user *entity.User) (*entity.User, error)
	Read() ([]*entity.User, error)
	ReadByID(userID uint) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(userID uint) error
	Login(email, password string) (string, error)
}

package user_usecase

import "github.com/celpung/gocleanarch/entity"

type UserUsecaseInterface interface {
	Create(user *entity.User) (*entity.UserHttpResponse, error)
	Read() ([]*entity.UserHttpResponse, error)
	ReadByID(userID uint) (*entity.UserHttpResponse, error)
	Update(user *entity.UserUpdate) (*entity.UserHttpResponse, error)
	Delete(userID uint) error
	Login(email, password string) (string, error)
}

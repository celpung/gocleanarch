package user_usecase

import user_entity "github.com/celpung/gocleanarch/domain/user/entity"

type UserUsecaseInterface interface {
	Create(user *user_entity.User) (*user_entity.UserHttpResponse, error)
	Read() ([]*user_entity.UserHttpResponse, error)
	ReadByID(userID uint) (*user_entity.UserHttpResponse, error)
	Update(user *user_entity.UserUpdate) (*user_entity.UserHttpResponse, error)
	Delete(userID uint) error
	Login(email, password string) (string, error)
}

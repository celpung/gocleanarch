package usecase

import user_entity "github.com/celpung/gocleanarch/domain/user/entity"

type UserUsecaseInterface interface {
	Create(user *user_entity.User) (*user_entity.User, error)
	Read() ([]*user_entity.User, error)
	ReadByID(userID uint) (*user_entity.User, error)
	Update(user *user_entity.User) (*user_entity.User, error)
	SoftDelete(userID uint) error
	Login(email, password string) (string, error)
}

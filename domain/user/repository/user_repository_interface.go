package user_repository

import user_entity "github.com/celpung/gocleanarch/domain/user/entity"

type UserRepositoryInterface interface {
	Create(user *user_entity.User) (*user_entity.User, error)
	Read() ([]*user_entity.User, error)
	ReadByID(userID uint) (*user_entity.User, error)
	ReadByEmail(email string, isLogin bool) (*user_entity.User, error)
	Update(user *user_entity.User) (*user_entity.User, error)
	Delete(userID uint) error
}

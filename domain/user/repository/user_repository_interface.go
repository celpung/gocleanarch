package user_repository

import "github.com/celpung/gocleanarch/domain/user/entity"

type UserRepositoryInterface interface {
	Create(user *entity.User) (*entity.User, error)
	Read() ([]*entity.User, error)
	ReadByID(userID uint) (*entity.User, error)
	ReadByEmail(email string, isLogin bool) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(userID uint) error
}

package usecase

import "github.com/celpung/gocleanarch/application/user/domain/entity"

type UserUsecase interface {
	Create(user *entity.User) (*entity.User, error)
	Read() ([]*entity.User, error)
	ReadByID(userID string) (*entity.User, error)
	// Update(user *entity.User) (*entity.User, error)
	Update(payload *entity.UpdateUserPayload) (*entity.User, error)
	SoftDelete(userID string) error
	Login(email, password string) (string, error)
}

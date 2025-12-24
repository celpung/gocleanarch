package port

import "github.com/celpung/gocleanarch/app/modules/user/domain/entity"

type UserUsecase interface {
	CreateUser(payload *entity.User) error
}

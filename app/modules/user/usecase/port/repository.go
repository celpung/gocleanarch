package port

import "github.com/celpung/gocleanarch/app/modules/user/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
}

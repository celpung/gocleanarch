package usecase

import (
	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/repository"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) CreateUser(user *entity.User) error {
	// do some logic here
	return uc.UserRepository.Create(user)
}

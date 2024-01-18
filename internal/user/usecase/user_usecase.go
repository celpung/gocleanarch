package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
	"github.com/celpung/gocleanarch/infrastructure"
)

type UserUsecase struct {
	UserRepository domain.UserRepository
	PasswordSrv    *infrastructure.PasswordService
	Jwtsrv         *infrastructure.JwtService
}

func NewUserUseCase(repository domain.UserRepository, passwordSrv *infrastructure.PasswordService, jwtSrv *infrastructure.JwtService) *UserUsecase {
	return &UserUsecase{
		UserRepository: repository,
		PasswordSrv:    passwordSrv,
		Jwtsrv:         jwtSrv,
	}
}

func (uc *UserUsecase) CreateUser(user *domain.User) error {
	hashedPassword, err := uc.PasswordSrv.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.UserRepository.Create(user)
}

// Read implements the Read method of the domain.UserUsecase interface.
func (uc *UserUsecase) Read(ctx context.Context) ([]domain.User, error) {
	return uc.UserRepository.Read(ctx)
}

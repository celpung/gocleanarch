package usecaseimplementation

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
	"github.com/celpung/gocleanarch/internal/user/repository"
	"github.com/celpung/gocleanarch/internal/user/usecase"
	"github.com/celpung/gocleanarch/services"
)

type UserUsecase struct {
	UserRepository repository.UserRepository
	PasswordSrv    *services.PasswordService
	Jwtsrv         *services.JwtService
}

func NewUserUseCase(repository repository.UserRepository, passwordSrv *services.PasswordService, jwtSrv *services.JwtService) usecase.UserUsecase {
	return &UserUsecase{
		UserRepository: repository,
		PasswordSrv:    passwordSrv,
		Jwtsrv:         jwtSrv,
	}
}

// CreateUser implements usecase.UserUsecase.
func (uc *UserUsecase) CreateUser(user *domain.User) error {
	hashedPassword, err := uc.PasswordSrv.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.UserRepository.Create(user)
}

// Read implements usecase.UserUsecase.
func (uc *UserUsecase) Read(ctx context.Context) ([]domain.User, error) {
	return uc.UserRepository.Read(ctx)
}

package usecase

import (
	"github.com/celpung/gocleanarch/infrastructure"
	"github.com/celpung/gocleanarch/internal/entity"
	userRepository "github.com/celpung/gocleanarch/internal/repository/user"
)

type UserUseCase struct {
	UserRepository userRepository.UserRepository
	PasswordSrv    *infrastructure.PasswordService
	Jwtsrv         *infrastructure.JwtService
}

func NewUserUseCase(userRepository userRepository.UserRepository, passwordSrv *infrastructure.PasswordService, jwtSrv *infrastructure.JwtService) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
		PasswordSrv:    passwordSrv,
		Jwtsrv:         jwtSrv,
	}
}

func (uc *UserUseCase) CreateUser(user *entity.User) error {
	hashedPassword, err := uc.PasswordSrv.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.UserRepository.Create(user)
}

func (uc *UserUseCase) Read() ([]entity.User, error) {
	return uc.UserRepository.Read()
}

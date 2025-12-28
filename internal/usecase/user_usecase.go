package usecase

import (
	"context"
	"errors"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

type UserUsecase struct {
	repo           repository.UserRepository
	passwordHasher dependencies.PasswordHasher
	jwtGenerator   dependencies.JwtGenerator
}

func (u *UserUsecase) Register(ctx context.Context, user entity.User) error {
	hash, err := u.passwordHasher.Hash(user.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	user.Password = hash

	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := u.passwordHasher.Compare(usr.Password, password); err != nil {
		return "", errors.New("password not match")
	}

	token, err := u.jwtGenerator.Generate(usr.ID, usr.Email, usr.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserUsecase(
	repo repository.UserRepository,
	passwordHasher dependencies.PasswordHasher,
	jwtGenerator dependencies.JwtGenerator) usecase.UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

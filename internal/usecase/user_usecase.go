package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type UserUsecase struct {
	repo           repository.UserRepository
	idGenerator    dependencies.IDGenerator
	passwordHasher dependencies.PasswordHasher
	jwtGenerator   dependencies.JwtGenerator
	typograph      dependencies.TypoGraph
}

func (u *UserUsecase) Register(ctx context.Context, user entity.User) error {
	hash, err := u.passwordHasher.Hash(user.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	uuid, err := u.idGenerator.NewID()
	if err != nil {
		return errors.New("failed to generate ID")
	}

	user.ID = uuid
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

func (u *UserUsecase) ChangePassword(ctx context.Context, email string, password string) error {
	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	hash, err := u.passwordHasher.Hash(usr.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	updates := make(map[string]any)
	if password != "" {
		updates["password"] = hash
	}

	return u.repo.Update(ctx, usr.ID, updates)
}

func (u *UserUsecase) UserLists(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	offset := (page - 1) * limit
	return u.repo.Lists(ctx, offset, limit)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, input *entity.User) error {
	updates := make(map[string]any)

	if input.Name != "" {
		name := strings.TrimSpace(typograph.ToTitleCase(input.Name))
		updates["name"] = name
	}

	if input.Email != "" {
		email := strings.TrimSpace(input.Email)
		updates["email"] = email
	}

	if input.Role != "" {
		role := strings.TrimSpace(input.Role)
		updates["role"] = role
	}

	return u.repo.Update(ctx, id, updates)
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func NewUserUsecase(
	repo repository.UserRepository,
	passwordHasher dependencies.PasswordHasher,
	jwtGenerator dependencies.JwtGenerator) usecase.UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

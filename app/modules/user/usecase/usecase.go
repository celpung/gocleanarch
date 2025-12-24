package usecase

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/celpung/gocleanarch/app/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/app/modules/user/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type UserUsecaseStruct struct {
	repo port.UserRepository
}

func (u *UserUsecaseStruct) CreateUser(payload *entity.User) error {
	if payload == nil {
		return errors.New("payload is required")
	}

	payload.Name = strings.TrimSpace(typograph.ToTitleCase(payload.Name))
	payload.Email = strings.TrimSpace(strings.ToLower(payload.Email))

	if payload.Name == "" {
		return errors.New("name is required")
	}
	if payload.Email == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(payload.Password) == "" {
		return errors.New("password is required")
	}

	if payload.ID == "" {
		payload.ID = uuid.NewString()
	}

	hashed, err := hashPassword(payload.Password)
	if err != nil {
		return err
	}
	
	payload.Password = hashed
	payload.IsActive = true

	return u.repo.CreateUser(payload)
}

func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewUserUsecase(repo port.UserRepository) port.UserUsecase {
	return &UserUsecaseStruct{repo: repo}
}

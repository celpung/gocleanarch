package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/celpung/gocleanarch/internal/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/internal/modules/user/usecase/dto"
	"github.com/celpung/gocleanarch/internal/modules/user/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type UserUsecaseStruct struct {
	repo   port.UserRepository
	hasher port.PasswordHasher
	idGen  port.IDGenerator
}

func (u *UserUsecaseStruct) CreateUser(ctx context.Context, payload *entity.User) error {
	if payload == nil {
		return port.ValidationError{Message: "payload is required"}
	}

	payload.Name = strings.TrimSpace(typograph.ToTitleCase(payload.Name))
	payload.Email = strings.TrimSpace(strings.ToLower(payload.Email))

	if payload.Name == "" {
		return port.ValidationError{Field: "name", Message: "is required"}
	}
	if payload.Email == "" {
		return port.ValidationError{Field: "email", Message: "is required"}
	}
	if strings.TrimSpace(payload.Password) == "" {
		return port.ValidationError{Field: "password", Message: "is required"}
	}

	if payload.ID == "" {
		id, err := u.idGen.NewID()
		if err != nil {
			return err
		}
		payload.ID = id
	}

	hashed, err := u.hasher.Hash(payload.Password)
	if err != nil {
		return err
	}

	payload.Password = hashed
	payload.IsActive = true

	return u.repo.CreateUser(ctx, payload)
}

func NewUserUsecase(repo port.UserRepository, hasher port.PasswordHasher, idGen port.IDGenerator) port.UserUsecase {
	return &UserUsecaseStruct{
		repo:   repo,
		hasher: hasher,
		idGen:  idGen,
	}
}

func (u *UserUsecaseStruct) Authenticate(ctx context.Context, email, password string) (*entity.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" {
		return nil, port.ValidationError{Field: "email", Message: "is required"}
	}
	if password == "" {
		return nil, port.ValidationError{Field: "password", Message: "is required"}
	}

	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return nil, port.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := u.hasher.Compare(user.Password, password); err != nil {
		return nil, port.ErrInvalidCredentials
	}

	return user, nil
}

func (u *UserUsecaseStruct) ListUsers(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	if page < 1 {
		return nil, 0, port.ValidationError{Field: "page", Message: "must be >= 1"}
	}
	if limit < 1 || limit > 100 {
		return nil, 0, port.ValidationError{Field: "limit", Message: "must be between 1 and 100"}
	}

	offset := (page - 1) * limit
	return u.repo.ListUsers(ctx, offset, limit)
}

func (u *UserUsecaseStruct) UpdateUser(ctx context.Context, id string, payload dto.UpdateUserInput) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}

	updates := make(map[string]any)

	if payload.Name != nil {
		name := strings.TrimSpace(typograph.ToTitleCase(*payload.Name))
		if name == "" {
			return port.ValidationError{Field: "name", Message: "is required"}
		}
		updates["name"] = name
	}

	if payload.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*payload.Email))
		if email == "" {
			return port.ValidationError{Field: "email", Message: "is required"}
		}
		updates["email"] = email
	}

	if payload.Password != nil {
		password := strings.TrimSpace(*payload.Password)
		if password == "" {
			return port.ValidationError{Field: "password", Message: "is required"}
		}
		hashed, err := u.hasher.Hash(password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}

	if len(updates) == 0 {
		return port.ErrNoFieldsToUpdate
	}

	return u.repo.UpdateUser(ctx, id, updates)
}

func (u *UserUsecaseStruct) DeleteUser(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}
	return u.repo.DeleteUser(ctx, id)
}

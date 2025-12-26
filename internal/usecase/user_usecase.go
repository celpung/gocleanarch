package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type UserUsecase struct {
	repo        port.UserRepository
	companyRepo port.CompanyRepository
	hasher      port.PasswordHasher
	idGen       port.IDGenerator
}

func NewUserUsecase(
	repo port.UserRepository,
	companyRepo port.CompanyRepository,
	hasher port.PasswordHasher,
	idGen port.IDGenerator,
) port.UserUsecase {
	return &UserUsecase{
		repo:        repo,
		companyRepo: companyRepo,
		hasher:      hasher,
		idGen:       idGen,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, input port.CreateUserInput) (*entity.User, error) {
	name := strings.TrimSpace(typograph.ToTitleCase(input.Name))
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	if name == "" {
		return nil, port.ValidationError{Field: "name", Message: "is required"}
	}
	if email == "" {
		return nil, port.ValidationError{Field: "email", Message: "is required"}
	}
	if password == "" {
		return nil, port.ValidationError{Field: "password", Message: "is required"}
	}

	companyID := ""
	if input.CompanyID != nil {
		companyID = strings.TrimSpace(*input.CompanyID)
		if companyID == "" {
			return nil, port.ValidationError{Field: "company_id", Message: "is required"}
		}

		if u.companyRepo != nil {
			if _, err := u.companyRepo.GetCompany(ctx, companyID); err != nil {
				return nil, err
			}
		}
	}

	id, err := u.idGen.NewID()
	if err != nil {
		return nil, err
	}

	hashed, err := u.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  hashed,
		CompanyID: companyID,
		IsActive:  true,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, port.ValidationError{Field: "id", Message: "is required"}
	}
	return u.repo.GetUser(ctx, id)
}

func (u *UserUsecase) ListUsers(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	if page < 1 {
		return nil, 0, port.ValidationError{Field: "page", Message: "must be >= 1"}
	}
	if limit < 1 || limit > 100 {
		return nil, 0, port.ValidationError{Field: "limit", Message: "must be between 1 and 100"}
	}

	offset := (page - 1) * limit
	return u.repo.ListUsers(ctx, offset, limit)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, input port.UpdateUserInput) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}

	updates := make(map[string]any)

	if input.Name != nil {
		name := strings.TrimSpace(typograph.ToTitleCase(*input.Name))
		if name == "" {
			return port.ValidationError{Field: "name", Message: "is required"}
		}
		updates["name"] = name
	}

	if input.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*input.Email))
		if email == "" {
			return port.ValidationError{Field: "email", Message: "is required"}
		}
		updates["email"] = email
	}

	if input.Password != nil {
		password := strings.TrimSpace(*input.Password)
		if password == "" {
			return port.ValidationError{Field: "password", Message: "is required"}
		}
		hashed, err := u.hasher.Hash(password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}

	if input.CompanyID != nil {
		companyID := strings.TrimSpace(*input.CompanyID)
		if companyID == "" {
			return port.ValidationError{Field: "company_id", Message: "is required"}
		}

		if u.companyRepo != nil {
			if _, err := u.companyRepo.GetCompany(ctx, companyID); err != nil {
				return err
			}
		}
		updates["company_id"] = companyID
	}

	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) == 0 {
		return port.ErrNoFieldsToUpdate
	}

	return u.repo.UpdateUser(ctx, id, updates)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}
	return u.repo.DeleteUser(ctx, id)
}

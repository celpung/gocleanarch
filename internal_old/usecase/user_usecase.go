package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal_old/entity"
	"github.com/celpung/gocleanarch/internal_old/usecase/port"
	porterrors "github.com/celpung/gocleanarch/internal_old/usecase/port/errors"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type UserUsecase struct {
	repo        port.UserRepository
	companyRepo port.CompanyRepository
	hasher      port.PasswordHasher
	idGen       port.IDGenerator
	JwtGen      port.JwtGenerator
}

func (u *UserUsecase) CreateUser(ctx context.Context, input port.CreateUserInput) (*entity.User, error) {
	name := strings.TrimSpace(typograph.ToTitleCase(input.Name))
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	companyID := ""
	if input.CompanyID != nil {
		companyID = strings.TrimSpace(*input.CompanyID)

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
		Role:      "user",
		CompanyID: companyID,
		IsActive:  true,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.GetUser(ctx, strings.TrimSpace(id))
}

func (u *UserUsecase) ListUsers(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	offset := (page - 1) * limit
	return u.repo.ListUsers(ctx, offset, limit)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, input port.UpdateUserInput) error {
	id = strings.TrimSpace(id)

	updates := make(map[string]any)

	if input.Name != nil {
		name := strings.TrimSpace(typograph.ToTitleCase(*input.Name))
		updates["name"] = name
	}

	if input.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*input.Email))
		updates["email"] = email
	}

	if input.Password != nil {
		password := strings.TrimSpace(*input.Password)
		hashed, err := u.hasher.Hash(password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}

	if input.CompanyID != nil {
		companyID := strings.TrimSpace(*input.CompanyID)

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
		return porterrors.ErrNoFieldsToUpdate
	}

	return u.repo.UpdateUser(ctx, id, updates)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, strings.TrimSpace(id))
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", porterrors.ErrUserNotFound
	}

	if err := u.hasher.Compare(usr.Password, password); err != nil {
		return "", porterrors.ErrWrongPassword
	}

	token, err := u.JwtGen.Generate(usr.ID, usr.Email, usr.Role)
	if err != nil {
		return "", porterrors.ErrJwtFailure
	}

	return token, nil
}

func NewUserUsecase(
	repo port.UserRepository,
	companyRepo port.CompanyRepository,
	hasher port.PasswordHasher,
	idGen port.IDGenerator,
	JwtGen port.JwtGenerator,
) port.UserUsecase {
	return &UserUsecase{
		repo:        repo,
		companyRepo: companyRepo,
		hasher:      hasher,
		idGen:       idGen,
		JwtGen:      JwtGen,
	}
}

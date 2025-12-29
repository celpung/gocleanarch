package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
	"github.com/celpung/gocleanarch/internal/domain/value"
	"github.com/celpung/gocleanarch/internal/usecase/dto"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	usecase_port "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

type UserUsecase struct {
	repo           repository.UserRepository
	idGenerator    dependencies.IDGenerator
	passwordHasher dependencies.PasswordHasher
	jwtGenerator   dependencies.JwtGenerator
	typograph      dependencies.TypoGraph
}

func NewUserUsecase(
	repo repository.UserRepository,
	idGenerator dependencies.IDGenerator,
	passwordHasher dependencies.PasswordHasher,
	jwtGenerator dependencies.JwtGenerator,
	typograph dependencies.TypoGraph,
) usecase_port.UserUsecase {
	return &UserUsecase{
		repo:           repo,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
		jwtGenerator:   jwtGenerator,
		typograph:      typograph,
	}
}

func (u *UserUsecase) Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error) {
	name := strings.TrimSpace(u.typograph.ToTitleCase(input.Name))
	if name == "" {
		return entity.User{}, apperrors.ErrNameRequired
	}

	email, err := value.NewEmail(input.Email)
	if err != nil {
		return entity.User{}, err
	}

	role, err := value.NewRole(input.Role)
	if err != nil {
		return entity.User{}, err
	}

	password, err := value.NewPassword(input.Password)
	if err != nil {
		return entity.User{}, err
	}

	hash, err := u.passwordHasher.Hash(password.String())
	if err != nil {
		return entity.User{}, apperrors.ErrPasswordHash
	}

	id, err := u.idGenerator.NewID()
	if err != nil {
		return entity.User{}, apperrors.ErrIDGeneration
	}

	user := entity.User{
		ID:       id,
		Name:     name,
		Email:    email.String(),
		Role:     role.String(),
		Password: hash,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return entity.User{}, err
	}

	created, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return *created, nil
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	emailVO, err := value.NewEmail(email)
	if err != nil {
		return "", err
	}

	passwordVO, err := value.NewPassword(password)
	if err != nil {
		return "", err
	}

	usr, err := u.repo.FindByEmail(ctx, emailVO.String())
	if err != nil {
		return "", err
	}

	if err := u.passwordHasher.Compare(usr.Password, passwordVO.String()); err != nil {
		return "", apperrors.ErrPasswordMismatch
	}

	token, err := u.jwtGenerator.Generate(usr.ID, usr.Email, usr.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) ChangePassword(ctx context.Context, userID string, password string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return apperrors.ErrUserIDRequired
	}

	passwordVO, err := value.NewPassword(password)
	if err != nil {
		return err
	}

	hash, err := u.passwordHasher.Hash(passwordVO.String())
	if err != nil {
		return apperrors.ErrPasswordHash
	}

	return u.repo.Update(ctx, userID, &entity.UpdateUser{
		Password: &hash,
	})
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (entity.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.User{}, apperrors.ErrUserIDRequired
	}

	usr, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return *usr, nil
}

func (u *UserUsecase) List(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	users, total, err := u.repo.Lists(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UserUsecase) Update(ctx context.Context, id string, updates dto.UpdateUserInput) (entity.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.User{}, apperrors.ErrUserIDRequired
	}

	var normalized entity.UpdateUser

	if updates.Name != nil {
		name := strings.TrimSpace(u.typograph.ToTitleCase(*updates.Name))
		if name != "" {
			normalized.Name = &name
		}
	}

	if updates.Email != nil {
		email, err := value.NewEmail(*updates.Email)
		if err != nil {
			return entity.User{}, err
		}
		emailStr := email.String()
		normalized.Email = &emailStr
	}

	if updates.Role != nil {
		role, err := value.NewRole(*updates.Role)
		if err != nil {
			return entity.User{}, err
		}
		roleStr := role.String()
		normalized.Role = &roleStr
	}

	if updates.Password != nil {
		password, err := value.NewPassword(*updates.Password)
		if err != nil {
			return entity.User{}, err
		}
		hash, err := u.passwordHasher.Hash(password.String())
		if err != nil {
			return entity.User{}, apperrors.ErrPasswordHash
		}
		normalized.Password = &hash
	}

	if normalized.Name == nil && normalized.Email == nil && normalized.Role == nil && normalized.Password == nil {
		return entity.User{}, apperrors.ErrNoChanges
	}

	if err := u.repo.Update(ctx, id, &normalized); err != nil {
		return entity.User{}, err
	}

	updated, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return *updated, nil
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperrors.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
}

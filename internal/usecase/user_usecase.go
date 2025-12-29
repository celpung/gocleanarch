package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
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

func (u *UserUsecase) Create(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	req.Name = strings.TrimSpace(u.typograph.ToTitleCase(req.Name))
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(req.Role)
	req.Password = strings.TrimSpace(req.Password)

	switch {
	case req.Name == "":
		return dto.UserResponse{}, apperrors.ErrNameRequired
	case req.Email == "":
		return dto.UserResponse{}, apperrors.ErrEmailRequired
	case req.Role == "":
		return dto.UserResponse{}, apperrors.ErrRoleRequired
	case req.Password == "":
		return dto.UserResponse{}, apperrors.ErrPasswordRequired
	}

	hash, err := u.passwordHasher.Hash(req.Password)
	if err != nil {
		return dto.UserResponse{}, apperrors.ErrPasswordHash
	}

	id, err := u.idGenerator.NewID()
	if err != nil {
		return dto.UserResponse{}, apperrors.ErrIDGeneration
	}

	user := entity.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: hash,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return dto.UserResponse{}, err
	}

	created, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(*created), nil
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" {
		return "", apperrors.ErrEmailRequired
	}
	if password == "" {
		return "", apperrors.ErrPasswordRequired
	}

	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := u.passwordHasher.Compare(usr.Password, password); err != nil {
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
	password = strings.TrimSpace(password)

	if userID == "" {
		return apperrors.ErrUserIDRequired
	}
	if password == "" {
		return apperrors.ErrPasswordRequired
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return apperrors.ErrPasswordHash
	}

	return u.repo.Update(ctx, userID, &entity.UpdateUser{
		Password: &hash,
	})
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (dto.UserResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return dto.UserResponse{}, apperrors.ErrUserIDRequired
	}

	usr, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(*usr), nil
}

func (u *UserUsecase) List(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error) {
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

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, toUserResponse(user))
	}

	return responses, total, nil
}

func (u *UserUsecase) Update(ctx context.Context, id string, updates dto.UpdateUserRequest) (dto.UserResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return dto.UserResponse{}, apperrors.ErrUserIDRequired
	}

	var normalized entity.UpdateUser

	if updates.Name != nil {
		name := strings.TrimSpace(u.typograph.ToTitleCase(*updates.Name))
		if name != "" {
			normalized.Name = &name
		}
	}

	if updates.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*updates.Email))
		if email != "" {
			normalized.Email = &email
		}
	}

	if updates.Role != nil {
		role := strings.TrimSpace(*updates.Role)
		if role != "" {
			normalized.Role = &role
		}
	}

	if updates.Password != nil {
		password := strings.TrimSpace(*updates.Password)
		if password == "" {
			return dto.UserResponse{}, apperrors.ErrPasswordRequired
		}
		hash, err := u.passwordHasher.Hash(password)
		if err != nil {
			return dto.UserResponse{}, apperrors.ErrPasswordHash
		}
		normalized.Password = &hash
	}

	if normalized.Name == nil && normalized.Email == nil && normalized.Role == nil && normalized.Password == nil {
		return dto.UserResponse{}, apperrors.ErrNoChanges
	}

	if err := u.repo.Update(ctx, id, &normalized); err != nil {
		return dto.UserResponse{}, err
	}

	updated, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(*updated), nil
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return apperrors.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
}

func toUserResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

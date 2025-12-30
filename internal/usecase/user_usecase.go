package usecase

import (
	"context"
	"net/mail"
	"strings"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	errs "github.com/celpung/gocleanarch/internal/domain/errors"
	"github.com/celpung/gocleanarch/internal/usecase/port"
)

type UserUsecase struct {
	repo port.UserRepository
	dep  port.Dependencies
}

func (u *UserUsecase) Create(ctx context.Context, input entity.User) (entity.User, error) {
	name := strings.TrimSpace(u.dep.ToTitleCase(input.Name))
	if name == "" {
		return entity.User{}, errs.ErrNameRequired
	}

	email, err := normalizeEmail(input.Email)
	if err != nil {
		return entity.User{}, err
	}

	role, err := normalizeRole(input.Role)
	if err != nil {
		return entity.User{}, err
	}

	password, err := normalizePassword(input.Password)
	if err != nil {
		return entity.User{}, err
	}

	hash, err := u.dep.HashPassword(password)
	if err != nil {
		return entity.User{}, errs.ErrPasswordHash
	}

	id, err := u.dep.NewID()
	if err != nil {
		return entity.User{}, errs.ErrIDGeneration
	}

	user := entity.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Role:     role,
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
	emailNormalized, err := normalizeEmail(email)
	if err != nil {
		return "", err
	}

	passwordNormalized, err := normalizePassword(password)
	if err != nil {
		return "", err
	}

	usr, err := u.repo.FindByEmail(ctx, emailNormalized)
	if err != nil {
		return "", err
	}

	if err := u.dep.ComparePassword(usr.Password, passwordNormalized); err != nil {
		return "", errs.ErrPasswordMismatch
	}

	token, err := u.dep.GenerateToken(usr.ID, usr.Email, usr.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) ChangePassword(ctx context.Context, userID string, password string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errs.ErrUserIDRequired
	}

	passwordNormalized, err := normalizePassword(password)
	if err != nil {
		return err
	}

	hash, err := u.dep.HashPassword(passwordNormalized)
	if err != nil {
		return errs.ErrPasswordHash
	}

	return u.repo.Update(ctx, userID, &entity.UpdateUser{
		Password: &hash,
	})
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (entity.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.User{}, errs.ErrUserIDRequired
	}

	usr, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return *usr, nil
}

func (u *UserUsecase) List(ctx context.Context, page int, limit int) ([]entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.repo.Lists(ctx, offset, limit)
}

func (u *UserUsecase) Update(ctx context.Context, id string, updates entity.UpdateUser) (entity.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.User{}, errs.ErrUserIDRequired
	}

	var normalized entity.UpdateUser

	if updates.Name != nil {
		name := strings.TrimSpace(*updates.Name)
		if name != "" {
			normalized.Name = &name
		}
	}

	if updates.Email != nil {
		email, err := normalizeEmail(*updates.Email)
		if err != nil {
			return entity.User{}, err
		}
		normalized.Email = &email
	}

	if updates.Role != nil {
		role, err := normalizeRole(*updates.Role)
		if err != nil {
			return entity.User{}, err
		}
		normalized.Role = &role
	}

	if updates.Password != nil {
		password, err := normalizePassword(*updates.Password)
		if err != nil {
			return entity.User{}, err
		}
		hash, err := u.dep.HashPassword(password)
		if err != nil {
			return entity.User{}, errs.ErrPasswordHash
		}
		normalized.Password = &hash
	}

	if normalized.Name == nil && normalized.Email == nil && normalized.Role == nil && normalized.Password == nil {
		return entity.User{}, errs.ErrNoChanges
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
		return errs.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
}

func NewUserUsecase(repo port.UserRepository,
	dep port.Dependencies,
) port.UserUsecase {
	return &UserUsecase{
		repo: repo,
		dep:  dep,
	}
}

func normalizeEmail(raw string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	if normalized == "" {
		return "", errs.ErrEmailRequired
	}
	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", errs.ErrInvalidEmail
	}
	return normalized, nil
}

func normalizeRole(raw string) (string, error) {
	normalized := strings.TrimSpace(strings.ToUpper(raw))
	if normalized == "" {
		return "", errs.ErrRoleRequired
	}

	switch normalized {
	case "SUPER", "ADMIN", "USER":
		return normalized, nil
	default:
		return "", errs.ErrInvalidRole
	}
}

func normalizePassword(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", errs.ErrPasswordRequired
	}
	if len(normalized) < 8 {
		return "", errs.ErrWeakPassword
	}
	return normalized, nil
}

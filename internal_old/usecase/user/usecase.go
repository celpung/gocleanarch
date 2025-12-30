package user

import (
	"context"
	"net/mail"
	"strings"

	domain "github.com/celpung/gocleanarch/internal_old/domain/user"
)

type Service struct {
	repo           Repository
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
	jwtGenerator   JWTGenerator
	typograph      Typograph
}

func NewUsecase(
	repo Repository,
	idGenerator IDGenerator,
	passwordHasher PasswordHasher,
	jwtGenerator JWTGenerator,
	typograph Typograph,
) Usecase {
	return &Service{
		repo:           repo,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
		jwtGenerator:   jwtGenerator,
		typograph:      typograph,
	}
}

func (u *Service) Create(ctx context.Context, input CreateUserInput) (UserOutput, error) {
	name := strings.TrimSpace(u.typograph.ToTitleCase(input.Name))
	if name == "" {
		return UserOutput{}, domain.ErrNameRequired
	}

	email, err := newEmail(input.Email)
	if err != nil {
		return UserOutput{}, err
	}

	role, err := newRole(input.Role)
	if err != nil {
		return UserOutput{}, err
	}

	password, err := newPassword(input.Password)
	if err != nil {
		return UserOutput{}, err
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return UserOutput{}, domain.ErrPasswordHash
	}

	id, err := u.idGenerator.NewID()
	if err != nil {
		return UserOutput{}, domain.ErrIDGeneration
	}

	user := domain.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Role:     role,
		Password: hash,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return UserOutput{}, err
	}

	created, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return UserOutput{}, err
	}

	return toUserOutput(*created), nil
}

func (u *Service) Login(ctx context.Context, email string, password string) (string, error) {
	emailNormalized, err := newEmail(email)
	if err != nil {
		return "", err
	}

	passwordNormalized, err := newPassword(password)
	if err != nil {
		return "", err
	}

	usr, err := u.repo.FindByEmail(ctx, emailNormalized)
	if err != nil {
		return "", err
	}

	if err := u.passwordHasher.Compare(usr.Password, passwordNormalized); err != nil {
		return "", domain.ErrPasswordMismatch
	}

	token, err := u.jwtGenerator.Generate(usr.ID, usr.Email, usr.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *Service) ChangePassword(ctx context.Context, userID string, password string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return domain.ErrUserIDRequired
	}

	passwordNormalized, err := newPassword(password)
	if err != nil {
		return err
	}

	hash, err := u.passwordHasher.Hash(passwordNormalized)
	if err != nil {
		return domain.ErrPasswordHash
	}

	return u.repo.Update(ctx, userID, &domain.UpdateUser{
		Password: &hash,
	})
}

func (u *Service) GetByID(ctx context.Context, id string) (UserOutput, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return UserOutput{}, domain.ErrUserIDRequired
	}

	usr, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return UserOutput{}, err
	}

	return toUserOutput(*usr), nil
}

func (u *Service) List(ctx context.Context, page, limit int) ([]UserOutput, int64, error) {
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

	return toUserOutputs(users), total, nil
}

func (u *Service) Update(ctx context.Context, id string, updates UpdateUserInput) (UserOutput, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return UserOutput{}, domain.ErrUserIDRequired
	}

	var normalized domain.UpdateUser

	if updates.Name != nil {
		name := strings.TrimSpace(u.typograph.ToTitleCase(*updates.Name))
		if name != "" {
			normalized.Name = &name
		}
	}

	if updates.Email != nil {
		email, err := newEmail(*updates.Email)
		if err != nil {
			return UserOutput{}, err
		}
		normalized.Email = &email
	}

	if updates.Role != nil {
		role, err := newRole(*updates.Role)
		if err != nil {
			return UserOutput{}, err
		}
		normalized.Role = &role
	}

	if updates.Password != nil {
		password, err := newPassword(*updates.Password)
		if err != nil {
			return UserOutput{}, err
		}
		hash, err := u.passwordHasher.Hash(password)
		if err != nil {
			return UserOutput{}, domain.ErrPasswordHash
		}
		normalized.Password = &hash
	}

	if normalized.Name == nil && normalized.Email == nil && normalized.Role == nil && normalized.Password == nil {
		return UserOutput{}, domain.ErrNoChanges
	}

	if err := u.repo.Update(ctx, id, &normalized); err != nil {
		return UserOutput{}, err
	}

	updated, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return UserOutput{}, err
	}

	return toUserOutput(*updated), nil
}

func (u *Service) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
}

func toUserOutput(user domain.User) UserOutput {
	return UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toUserOutputs(users []domain.User) []UserOutput {
	result := make([]UserOutput, 0, len(users))
	for i := range users {
		result = append(result, toUserOutput(users[i]))
	}
	return result
}

const minPasswordLength = 8

func newEmail(raw string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	if normalized == "" {
		return "", domain.ErrEmailRequired
	}
	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", domain.ErrInvalidEmail
	}
	return normalized, nil
}

func newRole(raw string) (string, error) {
	normalized := strings.TrimSpace(strings.ToUpper(raw))
	if normalized == "" {
		return "", domain.ErrRoleRequired
	}
	switch normalized {
	case "SUPER", "ADMIN", "USER":
		return normalized, nil
	default:
		return "", domain.ErrInvalidRole
	}
}

func newPassword(raw string) (string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", domain.ErrPasswordRequired
	}
	if len(normalized) < minPasswordLength {
		return "", domain.ErrWeakPassword
	}
	return normalized, nil
}

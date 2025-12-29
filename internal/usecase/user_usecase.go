package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/entity"
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

func (u *UserUsecase) Register(ctx context.Context, user entity.User) error {
	// normalize & validate (minimal tapi penting)
	user.Name = strings.TrimSpace(u.typograph.ToTitleCase(user.Name))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Role = strings.TrimSpace(user.Role)
	user.Password = strings.TrimSpace(user.Password)

	if user.Name == "" {
		return entity.ErrNameRequired
	}
	if user.Email == "" {
		return entity.ErrEmailRequired
	}
	if user.Role == "" {
		return entity.ErrRoleRequired
	}
	if user.Password == "" {
		return entity.ErrPasswordRequired
	}

	hash, err := u.passwordHasher.Hash(user.Password)
	if err != nil {
		return entity.ErrPasswordHash
	}

	uuid, err := u.idGenerator.NewID()
	if err != nil {
		return entity.ErrIDGeneration
	}

	user.ID = uuid
	user.Password = hash

	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if email == "" {
		return "", entity.ErrEmailRequired
	}
	if password == "" {
		return "", entity.ErrPasswordRequired
	}

	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := u.passwordHasher.Compare(usr.Password, password); err != nil {
		return "", entity.ErrPasswordMismatch
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
		return entity.ErrUserIDRequired
	}
	if password == "" {
		return entity.ErrPasswordRequired
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return entity.ErrPasswordHash
	}

	return u.repo.Update(ctx, userID, &entity.UpdateUser{
		Password: &hash,
	})
}

func (u *UserUsecase) UserLists(ctx context.Context, page, limit int) ([]entity.User, int64, error) {
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

	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, input *entity.UpdateUser) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.ErrUserIDRequired
	}
	if input == nil {
		return entity.ErrInvalidInput
	}

	var updates entity.UpdateUser

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		name := strings.TrimSpace(u.typograph.ToTitleCase(*input.Name))
		updates.Name = &name
	}

	if input.Email != nil && strings.TrimSpace(*input.Email) != "" {
		email := strings.TrimSpace(strings.ToLower(*input.Email))
		updates.Email = &email
	}

	if input.Role != nil && strings.TrimSpace(*input.Role) != "" {
		role := strings.TrimSpace(*input.Role)
		updates.Role = &role
	}

	updates.Password = nil

	if updates.Name == nil && updates.Email == nil && updates.Role == nil {
		return entity.ErrNoChanges
	}

	return u.repo.Update(ctx, id, &updates)
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
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

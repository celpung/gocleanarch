package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	usecase_port "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	userdto "github.com/celpung/gocleanarch/internal/usecase/user/dto"
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

func (u *UserUsecase) Create(ctx context.Context, req userdto.CreateUserRequest) (userdto.UserResponse, error) {
	req.Name = strings.TrimSpace(u.typograph.ToTitleCase(req.Name))
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(req.Role)
	req.Password = strings.TrimSpace(req.Password)

	switch {
	case req.Name == "":
		return userdto.UserResponse{}, entity.ErrNameRequired
	case req.Email == "":
		return userdto.UserResponse{}, entity.ErrEmailRequired
	case req.Role == "":
		return userdto.UserResponse{}, entity.ErrRoleRequired
	case req.Password == "":
		return userdto.UserResponse{}, entity.ErrPasswordRequired
	}

	hash, err := u.passwordHasher.Hash(req.Password)
	if err != nil {
		return userdto.UserResponse{}, entity.ErrPasswordHash
	}

	id, err := u.idGenerator.NewID()
	if err != nil {
		return userdto.UserResponse{}, entity.ErrIDGeneration
	}

	user := entity.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: hash,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return userdto.UserResponse{}, err
	}

	created, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return userdto.UserResponse{}, err
	}

	return toUserResponse(*created), nil
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

func (u *UserUsecase) GetByID(ctx context.Context, id string) (userdto.UserResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return userdto.UserResponse{}, entity.ErrUserIDRequired
	}

	usr, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return userdto.UserResponse{}, err
	}

	return toUserResponse(*usr), nil
}

func (u *UserUsecase) List(ctx context.Context, page, limit int) ([]userdto.UserResponse, int64, error) {
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

	responses := make([]userdto.UserResponse, 0, len(users))
	for _, usr := range users {
		responses = append(responses, toUserResponse(usr))
	}

	return responses, total, nil
}

func (u *UserUsecase) Update(ctx context.Context, req userdto.UpdateUserRequest) (userdto.UserResponse, error) {
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		return userdto.UserResponse{}, entity.ErrUserIDRequired
	}

	var updates entity.UpdateUser

	if req.Name != nil {
		name := strings.TrimSpace(u.typograph.ToTitleCase(*req.Name))
		if name != "" {
			updates.Name = &name
		}
	}

	if req.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*req.Email))
		if email != "" {
			updates.Email = &email
		}
	}

	if req.Role != nil {
		role := strings.TrimSpace(*req.Role)
		if role != "" {
			updates.Role = &role
		}
	}

	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password == "" {
			return userdto.UserResponse{}, entity.ErrPasswordRequired
		}
		hash, err := u.passwordHasher.Hash(password)
		if err != nil {
			return userdto.UserResponse{}, entity.ErrPasswordHash
		}
		updates.Password = &hash
	}

	if updates.Name == nil && updates.Email == nil && updates.Role == nil && updates.Password == nil {
		return userdto.UserResponse{}, entity.ErrNoChanges
	}

	if err := u.repo.Update(ctx, req.ID, &updates); err != nil {
		return userdto.UserResponse{}, err
	}

	updated, err := u.repo.FindByID(ctx, req.ID)
	if err != nil {
		return userdto.UserResponse{}, err
	}

	return toUserResponse(*updated), nil
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return entity.ErrUserIDRequired
	}
	return u.repo.Delete(ctx, id)
}

func toUserResponse(user entity.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

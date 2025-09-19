package usecase_impl

import (
	"errors"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/repository"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/model"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
	"github.com/celpung/gocleanarch/infrastructure/typograph"
)

type UserUsecaseStruct struct {
	Repo            repository.UserRepository
	PasswordService *auth.PasswordService
	JWTService      *auth.JwtService
}

func (u *UserUsecaseStruct) Create(user *entity.User) (*entity.User, error) {
	hashed, err := u.PasswordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashed

	var m model.User
	if err := mapper.CopyTo(user, &m); err != nil {
		return nil, err
	}

	created, err := u.Repo.Create(&m)
	if err != nil {
		return nil, err
	}

	var out entity.User
	if err := mapper.CopyTo(created, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (u *UserUsecaseStruct) Update(payload *entity.UpdateUserPayload) (*entity.User, error) {
	if _, err := u.Repo.ReadByID(payload.ID); err != nil {
		return nil, err
	}

	changes := make(map[string]any)

	if payload.Name != nil {
		changes["name"] = *payload.Name
	}

	if payload.Email != nil {
		changes["email"] = *payload.Email
	}

	if payload.Password != nil {
		hashed, err := u.PasswordService.HashPassword(*payload.Password)
		if err != nil {
			return nil, err
		}
		changes["password"] = hashed
	}

	if payload.Active != nil {
		changes["active"] = *payload.Active
	}

	if payload.Role != nil {
		changes["role"] = *payload.Role
	}

	if len(changes) == 0 {
		cur, err := u.Repo.ReadByID(payload.ID)
		if err != nil {
			return nil, err
		}

		var out entity.User
		if err := mapper.CopyTo(cur, &out); err != nil {
			return nil, err
		}

		return &out, nil
	}

	updated, err := u.Repo.UpdateFields(payload.ID, changes)
	if err != nil {
		return nil, err
	}

	var res entity.User
	if err := mapper.CopyTo(updated, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (u *UserUsecaseStruct) SoftDelete(userID string) error {
	return u.Repo.SoftDelete(userID)
}

func (u *UserUsecaseStruct) Read(page, limit uint) ([]*entity.User, int64, error) {
	ms, total, err := u.Repo.Read(page, limit)
	if err != nil {
		return nil, 0, err
	}

	es := make([]*entity.User, len(ms))
	for i := range ms {
		es[i] = &entity.User{}
		_ = mapper.CopyTo(ms[i], es[i])
	}

	return es, total, nil
}

func (u *UserUsecaseStruct) ReadByID(userID string) (*entity.User, error) {
	m, err := u.Repo.ReadByID(userID)
	if err != nil {
		return nil, err
	}

	var out entity.User
	_ = mapper.CopyTo(m, &out)

	titleCased := typograph.ToTitleCase(out.Name)
	out.Name = titleCased

	return &out, nil
}

func (u *UserUsecaseStruct) Search(page, limit uint, keyword string) ([]*entity.User, int64, error) {
	ms, total, err := u.Repo.Search(page, limit, keyword)
	if err != nil {
		return nil, 0, err
	}

	es := make([]*entity.User, len(ms))
	for i := range ms {
		es[i] = &entity.User{}
		_ = mapper.CopyTo(ms[i], es[i])
	}

	return es, total, nil
}

func (u *UserUsecaseStruct) Login(email, password string) (string, error) {
	m, err := u.Repo.ReadByEmailPrivate(email)
	if err != nil {
		return "", err
	}

	if !m.Active {
		return "", errors.New("user not active")
	}

	if err := u.PasswordService.VerifyPassword(m.Password, password); err != nil {
		return "", errors.New("wrong password")
	}

	var e entity.User
	if err := mapper.CopyTo(m, &e); err != nil {
		return "", err
	}

	token, err := u.JWTService.JWTGenerator(e)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserUsecase(repo repository.UserRepository, passwordService *auth.PasswordService, jwtService *auth.JwtService) usecase.UserUsecase {
	return &UserUsecaseStruct{
		Repo:            repo,
		PasswordService: passwordService,
		JWTService:      jwtService,
	}
}

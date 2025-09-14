package usecase_impl

import (
	"errors"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/repository"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/model"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
)

type UserUsecaseStruct struct {
	Repo            repository.UserRepository
	PasswordService *auth.PasswordService
	JWTService      *auth.JwtService
}

func (u *UserUsecaseStruct) Create(user *entity.User) (*entity.User, error) {
	hashedPassword, err := u.PasswordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	var m model.User
	if err := mapper.CopyTo(user, &m); err != nil {
		return nil, err
	}

	usr, err := u.Repo.Create(&m)
	if err != nil {
		return nil, err
	}

	var res entity.User
	if err := mapper.CopyTo(usr, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (u *UserUsecaseStruct) Read() ([]*entity.User, error) {
	users, err := u.Repo.Read()
	if err != nil {
		return nil, err
	}

	res, err := mapper.MapStructList[model.User, entity.User](users)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserUsecaseStruct) ReadByID(userID string) (*entity.User, error) {
	user, err := u.Repo.ReadByID(userID)
	if err != nil {
		return nil, err
	}

	var res entity.User
	if err := mapper.CopyTo(user, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (u *UserUsecaseStruct) Update(payload *entity.UpdateUserPayload) (*entity.User, error) {
	_, err := u.Repo.ReadByID(payload.ID)
	if err != nil {
		return nil, err
	}

	changes := make(map[string]interface{})

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
		out.Password = ""
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

func (u *UserUsecaseStruct) Login(email, password string) (string, error) {
	user, err := u.Repo.ReadByEmailPrivate(email)
	if err != nil {
		return "", err
	}

	if !user.Active {
		return "", errors.New("user not active")
	}

	if err := u.PasswordService.VerifyPassword(user.Password, password); err != nil {
		return "", errors.New("wrong password")
	}

	var res entity.User
	if err := mapper.CopyTo(user, &res); err != nil {
		return "", err
	}

	token, err := u.JWTService.JWTGenerator(res)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserUsecase(repository repository.UserRepository, passwordService *auth.PasswordService, jwtService *auth.JwtService) usecase.UserUsecase {
	return &UserUsecaseStruct{
		Repo:            repository,
		PasswordService: passwordService,
		JWTService:      jwtService,
	}
}

package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
)

type UserUsecase struct {
	repo port.UserRepository
}

func (u *UserUsecase) Create(ctx context.Context, input entity.User) (entity.User, error) {
	panic("unimplemented")
}

func (u *UserUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	panic("unimplemented")
}

func (u *UserUsecase) ChangePassword(ctx context.Context, userID string, password string) error {
	panic("unimplemented")
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (entity.User, error) {
	panic("unimplemented")
}

func (u *UserUsecase) List(ctx context.Context, page int, limit int) ([]entity.User, int64, error) {
	panic("unimplemented")
}

func (u *UserUsecase) Update(ctx context.Context, id string, updates entity.UpdateUser) (entity.User, error) {
	panic("unimplemented")
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

func NewUserUsecase(repo port.UserRepository) port.UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

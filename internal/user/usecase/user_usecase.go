package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	Read(ctx context.Context) ([]domain.User, error)
}

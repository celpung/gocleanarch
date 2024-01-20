package repository

import (
	"context"

	"github.com/celpung/gocleanarch/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	Read(ctx context.Context) (user []domain.User, err error)
}

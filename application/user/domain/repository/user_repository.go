package repository

import (
	"github.com/celpung/gocleanarch/infrastructure/db/model"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	Read(page, limit uint) ([]*model.User, int64, error)
	ReadByID(userID string) (*model.User, error)
	ReadByEmailPublic(email string) (*model.User, error)
	ReadByEmailPrivate(email string) (*model.User, error)
	Search(page, limit uint, keyword string) ([]*model.User, int64, error)
	Update(user *model.User) (*model.User, error)
	UpdateFields(id string, fields map[string]any) (*model.User, error)
	SoftDelete(userID string) error
}

package repository

import (
	"github.com/celpung/gocleanarch/infrastructure/db/model"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	Read() ([]*model.User, error)
	ReadByID(userID string) (*model.User, error)
	ReadByEmailPublic(email string) (*model.User, error)
	ReadByEmailPrivate(email string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	UpdateFields(id string, fields map[string]interface{}) (*model.User, error)
	SoftDelete(userID string) error
}

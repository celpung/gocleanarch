package repository

import (
	"github.com/celpung/gocleanarch/app/infra/db/model"
	"github.com/celpung/gocleanarch/app/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/app/modules/user/usecase/port"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(user *entity.User) error {

	m := model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsActive: user.IsActive,
	}

	return r.db.Create(&m).Error
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{
		db: db,
	}
}

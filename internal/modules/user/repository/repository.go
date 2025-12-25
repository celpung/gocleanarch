package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/celpung/gocleanarch/internal/modules/user/domain/entity"
	"github.com/celpung/gocleanarch/internal/modules/user/usecase/port"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	m := model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsActive: user.IsActive,
	}

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		if isDuplicateEmailError(err) {
			return port.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func isDuplicateEmailError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		errors.Is(err, gorm.ErrDuplicatedKey)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	user := &entity.User{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Password: m.Password,
		IsActive: m.IsActive,
	}

	return user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	var models []model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, 0, len(models))
	for _, m := range models {
		users = append(users, entity.User{
			ID:       m.ID,
			Name:     m.Name,
			Email:    m.Email,
			Password: m.Password,
			IsActive: m.IsActive,
		})
	}

	return users, total, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if isDuplicateEmailError(result.Error) {
			return port.ErrEmailAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return port.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return port.ErrUserNotFound
	}
	return nil
}

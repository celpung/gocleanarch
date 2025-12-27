package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateEmailError(err) {
			return port.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	var (
		users []entity.User
		total int64
	)

	if err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Updates(fields)
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
	result := r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return port.ErrUserNotFound
	}

	return nil
}

func isDuplicateEmailError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || errors.Is(err, gorm.ErrDuplicatedKey)
}

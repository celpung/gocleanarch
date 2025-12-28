package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/celpung/gocleanarch/internal_old/configs/db/model"
	"github.com/celpung/gocleanarch/internal_old/entity"
	"github.com/celpung/gocleanarch/internal_old/usecase/port"
	porterrors "github.com/celpung/gocleanarch/internal_old/usecase/port/errors"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	record := toDBUser(*user)

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		if isDuplicateEmailError(err) {
			return porterrors.ErrEmailAlreadyExists
		}
		return err
	}

	// Propagate any defaults (e.g., role) back to entity.
	user.Role = record.Role

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*entity.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, porterrors.ErrUserNotFound
		}
		return nil, err
	}

	return toEntityUser(user), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, porterrors.ErrUserNotFound
		}
		return nil, err
	}

	return toEntityUser(user), nil
}

func (r *UserRepository) ListUsers(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	var (
		users []model.User
		total int64
	)

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entity.User, 0, len(users))
	for _, user := range users {
		result = append(result, *toEntityUser(user))
	}

	return result, total, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if isDuplicateEmailError(result.Error) {
			return porterrors.ErrEmailAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return porterrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return porterrors.ErrUserNotFound
	}

	return nil
}

func isDuplicateEmailError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || errors.Is(err, gorm.ErrDuplicatedKey)
}

func toEntityUser(user model.User) *entity.User {
	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		IsActive:  user.IsActive,
	}
}

func toDBUser(user entity.User) model.User {
	role := user.Role
	if role == "" {
		role = "user"
	}

	return model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      role,
		CompanyID: user.CompanyID,
		IsActive:  user.IsActive,
	}
}

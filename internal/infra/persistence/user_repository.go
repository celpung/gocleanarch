package persistence

import (
	"context"
	"errors"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	m := toModelUser(user)

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entity.ErrEmailExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) Lists(ctx context.Context, offset, limit int) ([]entity.User, int64, error) {
	var (
		users []model.User
		total int64
	)

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return toEntityUsers(users), total, nil
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*entity.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	e := toEntityUser(user)
	return &e, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrEmailNotFound
		}
		return nil, err
	}

	e := toEntityUser(user)
	return &e, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, input *entity.UpdateUser) error {
	if input == nil {
		return nil
	}

	updates := make(map[string]any)

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Role != nil {
		updates["role"] = *input.Role
	}
	if input.Password != nil {
		updates["password"] = *input.Password
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		if input.Email != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entity.ErrEmailExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Delete(&model.User{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// mappers
func toModelUser(user entity.User) model.User {
	return model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func toEntityUser(user model.User) entity.User {
	return entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func toEntityUsers(users []model.User) []entity.User {
	result := make([]entity.User, 0, len(users))
	for _, user := range users {
		result = append(result, toEntityUser(user))
	}
	return result
}

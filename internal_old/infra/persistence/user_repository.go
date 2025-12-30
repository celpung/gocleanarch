package persistence

import (
	"context"
	"errors"
	"strings"

	domain "github.com/celpung/gocleanarch/internal_old/domain/user"
	"github.com/celpung/gocleanarch/internal_old/infra/db/model"
	"github.com/celpung/gocleanarch/internal_old/usecase/user"
	"github.com/celpung/gocleanarch/pkg/mapper"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	m, err := toModelUser(user)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		if isDuplicateErr(err) {
			return domain.ErrEmailExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) Lists(ctx context.Context, offset, limit int) ([]domain.User, int64, error) {
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

	entities, err := toEntityUsers(users)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	e, err := toEntityUser(user)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmailNotFound
		}
		return nil, err
	}

	e, err := toEntityUser(user)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, input *domain.UpdateUser) error {
	if input == nil {
		return domain.ErrInvalidInput
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
		if input.Email != nil && isDuplicateErr(result.Error) {
			return domain.ErrEmailExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
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
		return domain.ErrUserNotFound
	}

	return nil
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

// mappers
func toModelUser(user domain.User) (model.User, error) {
	var m model.User
	if err := mapper.CopyTo(&user, &m); err != nil {
		return model.User{}, err
	}
	return m, nil
}

func toEntityUser(user model.User) (domain.User, error) {
	var e domain.User
	if err := mapper.CopyTo(&user, &e); err != nil {
		return domain.User{}, err
	}
	return e, nil
}

func toEntityUsers(users []model.User) ([]domain.User, error) {
	result := make([]domain.User, 0, len(users))
	for i := range users {
		e, err := toEntityUser(users[i])
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") ||
		strings.Contains(msg, "duplicate entry")
}

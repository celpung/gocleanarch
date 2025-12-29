package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
)

// UserRepository is an in-memory implementation for fast unit tests.
type UserRepository struct {
	mu    sync.Mutex
	store map[string]entity.User
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository() repository.UserRepository {
	return &UserRepository{store: make(map[string]entity.User)}
}

func (r *UserRepository) Create(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.store {
		if strings.EqualFold(u.Email, user.Email) {
			return entity.ErrEmailExists
		}
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	r.store[user.ID] = user
	return nil
}

func (r *UserRepository) Lists(_ context.Context, offset, limit int) ([]entity.User, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]entity.User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, cloneUser(u))
	}

	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i].Name) < strings.ToLower(users[j].Name)
	})

	total := int64(len(users))
	if offset > len(users) {
		return []entity.User{}, total, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], total, nil
}

func (r *UserRepository) FindByID(_ context.Context, userID string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.store[userID]
	if !ok {
		return nil, entity.ErrUserNotFound
	}
	cloned := cloneUser(user)
	return &cloned, nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.store {
		if strings.EqualFold(u.Email, email) {
			cloned := cloneUser(u)
			return &cloned, nil
		}
	}
	return nil, entity.ErrEmailNotFound
}

func (r *UserRepository) Update(_ context.Context, id string, input *entity.UpdateUser) error {
	if input == nil {
		return entity.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.store[id]
	if !ok {
		return entity.ErrUserNotFound
	}

	if input.Email != nil {
		for otherID, other := range r.store {
			if otherID != id && strings.EqualFold(other.Email, *input.Email) {
				return entity.ErrEmailExists
			}
		}
		user.Email = *input.Email
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Password != nil {
		user.Password = *input.Password
	}

	user.UpdatedAt = time.Now()
	r.store[id] = user

	return nil
}

func (r *UserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return entity.ErrUserNotFound
	}
	delete(r.store, id)
	return nil
}

func cloneUser(u entity.User) entity.User {
	return entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

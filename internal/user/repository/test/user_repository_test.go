package test

import (
	"context"
	"testing"

	"github.com/celpung/gocleanarch/domain"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository is a mock implementation of UserRepository for testing purposes.
type MockUserRepository struct {
	CreatedUsers []domain.User // Store created users for testing purposes
}

func (m *MockUserRepository) Create(user *domain.User) error {
	// Mock implementation: just append the user to the list
	m.CreatedUsers = append(m.CreatedUsers, *user)
	return nil
}

func (m *MockUserRepository) Read(ctx context.Context) ([]domain.User, error) {
	// Mock implementation: return the stored created users
	return m.CreatedUsers, nil
}

func TestCreateUser(t *testing.T) {
	repo := &MockUserRepository{} // Use the mock repository for testing
	user := &domain.User{}         // Create a sample user

	err := repo.Create(user)

	assert.NoError(t, err, "Expected no error when creating a user")
	assert.Len(t, repo.CreatedUsers, 1, "Expected one user to be created")
}

func TestReadUsers(t *testing.T) {
	repo := &MockUserRepository{} // Use the mock repository for testing
	ctx := context.Background()

	users, err := repo.Read(ctx)

	assert.NoError(t, err, "Expected no error when reading users")
	assert.Empty(t, users, "Expected empty list of users")
}
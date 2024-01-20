package test

import (
	"context"
	"testing"

	"github.com/celpung/gocleanarch/domain"
	"github.com/celpung/gocleanarch/internal/user/repository"
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

// MockUserUsecase is a mock implementation of UserUsecase for testing purposes.
type MockUserUsecase struct {
	Repo repository.UserRepository // Embed the UserRepository in the mock use case
}

func (m *MockUserUsecase) CreateUser(user *domain.User) error {
	// Delegate the call to the repository's Create method
	return m.Repo.Create(user)
}

func (m *MockUserUsecase) Read(ctx context.Context) ([]domain.User, error) {
	// Delegate the call to the repository's Read method
	return m.Repo.Read(ctx)
}

func TestCreateUser(t *testing.T) {
	repo := &MockUserRepository{}           // Use the mock repository for testing
	usecase := &MockUserUsecase{Repo: repo} // Inject the mock repository into the use case
	user := &domain.User{}                  // Create a sample user

	err := usecase.CreateUser(user)

	assert.NoError(t, err, "Expected no error when creating a user")
	assert.Len(t, repo.CreatedUsers, 1, "Expected one user to be created")
}

func TestReadUsers(t *testing.T) {
	repo := &MockUserRepository{}           // Use the mock repository for testing
	usecase := &MockUserUsecase{Repo: repo} // Inject the mock repository into the use case
	ctx := context.Background()

	users, err := usecase.Read(ctx)

	assert.NoError(t, err, "Expected no error when reading users")
	assert.Empty(t, users, "Expected empty list of users")
}

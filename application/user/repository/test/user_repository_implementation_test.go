package user_repository_implementation_test

import (
	"testing"

	user_repository_implementation "github.com/celpung/gocleanarch/application/user/repository"
	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	user_model "github.com/celpung/gocleanarch/infrastructure/db/model"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&user_model.User{})
	require.NoError(t, err)

	return db
}

func createUserEntity(name, email string) *user_entity.User {
	return &user_entity.User{
		Name:     name,
		Email:    email,
		Password: "password123",
		Active:   true,
		Role:     1,
	}
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	user := createUserEntity("Alice", "alice@example.com")

	saved, err := repo.Create(user)
	require.NoError(t, err)
	require.NotZero(t, saved.ID)
	require.Equal(t, "Alice", saved.Name)
}

func TestGetAllUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	usersToCreate := []*user_entity.User{
		createUserEntity("Maria", "maria@example.com"),
		createUserEntity("Bob", "bob@example.com"),
	}

	for _, user := range usersToCreate {
		_, err := repo.Create(user)
		require.NoError(t, err)
	}

	users, err := repo.Read()
	require.NoError(t, err)
	require.Len(t, users, len(usersToCreate))
	for i, user := range users {
		require.Equal(t, usersToCreate[i].Name, user.Name)
		require.Equal(t, usersToCreate[i].Email, user.Email)
	}
}

func TestReadByIDUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	user := createUserEntity("Charlie", "charlie@example.com")

	saved, err := repo.Create(user)
	require.NoError(t, err)

	usr, err := repo.ReadByID(1)
	require.NoError(t, err)
	require.Equal(t, saved.ID, usr.ID)
}

func TestReadByEmailUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	user := createUserEntity("Richard", "richard@example.com")

	saved, err := repo.Create(user)
	require.NoError(t, err)

	usr, err := repo.ReadByID(1)
	require.NoError(t, err)
	require.Equal(t, saved.Email, usr.Email)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	user := createUserEntity("Diana", "diana@example.com")

	saved, err := repo.Create(user)
	require.NoError(t, err)

	saved.Name = "Diana Updated"
	updated, err := repo.Update(saved)
	require.NoError(t, err)
	require.Equal(t, "Diana Updated", updated.Name)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository_implementation.NewUserRepository(db)

	user := createUserEntity("Eve", "eve@example.com")
	saved, err := repo.Create(user)
	require.NoError(t, err)

	err = repo.SoftDelete(saved.ID)
	require.NoError(t, err)
	_, err = repo.ReadByID(saved.ID)
	require.Error(t, err, "record not found")
	require.Equal(t, gorm.ErrRecordNotFound, err)
	
	// Ensure the user is not returned in Read
	users, err := repo.Read()
	require.NoError(t, err)
	require.NotContains(t, users, saved, "Deleted user should not be in the list")
	require.Len(t, users, 0, "No users should be present after deletion")
}

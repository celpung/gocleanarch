package user_repository_implementation_test

import (
	"testing"

	sqlite_configs "github.com/celpung/gocleanarch/configs/database/sqlite"
	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	user_repository_implementation "github.com/celpung/gocleanarch/domain/user/repository/implementation"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_create_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create a user entity for testing
	newUser := &user_entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	// Call the Create method
	createdUser, err := userRepository.Create(newUser)

	// Check if there is an error
	assert.NoError(t, err)
	// Check if the returned user ID is not zero
	assert.NotEqual(t, uint(0), createdUser.ID)
	// check other fields of the created user
	assert.Equal(t, newUser.Name, createdUser.Name)
	assert.Equal(t, newUser.Email, createdUser.Email)
	assert.Equal(t, newUser.Password, createdUser.Password)
	assert.Equal(t, newUser.Active, createdUser.Active)
	assert.Equal(t, newUser.Role, createdUser.Role)
}

func TestGetAllUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_read_all_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create some user entities for testing
	usersToCreate := []*user_entity.User{
		{Name: "User 1", Email: "user1@example.com", Password: "password1", Active: true, Role: 1},
		{Name: "User 2", Email: "user2@example.com", Password: "password2", Active: true, Role: 2},
		{Name: "User 3", Email: "user3@example.com", Password: "password3", Active: true, Role: 3},
	}

	// Create users in the database
	for _, user := range usersToCreate {
		_, err := userRepository.Create(user)
		assert.NoError(t, err)
	}

	// Call the Read method to get all users
	users, err := userRepository.Read()

	// Check if there is no error
	assert.NoError(t, err)
	// Check if the number of returned users matches the number of users created
	assert.Equal(t, len(usersToCreate), len(users))

	// Check individual properties of the returned users
	for i, user := range usersToCreate {
		assert.Equal(t, user.Name, users[i].Name)
		assert.Equal(t, user.Email, users[i].Email)
		assert.Equal(t, user.Password, users[i].Password)
		assert.Equal(t, user.Active, users[i].Active)
		assert.Equal(t, user.Role, users[i].Role)
	}
}

func TestReadByIdUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_read_by_id_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create a user entity for testing
	newUser := &user_entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	userRepository.Create(newUser)

	// Call the Read method to get user by id
	user, err := userRepository.ReadByID(1)

	// Check if there is no error
	assert.NoError(t, err)
	// check is id of checked by id is same with created user
	assert.Equal(t, user.ID, newUser.ID)
	assert.Equal(t, user.Name, newUser.Name)
	assert.Equal(t, user.Email, newUser.Email)
}

func TestReadByEmailUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_read_by_email_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create a user entity for testing
	newUser := &user_entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	userRepository.Create(newUser)

	// Call the Read method to get user by id
	user, err := userRepository.ReadByEmail(newUser.Email, false)

	// Check if there is no error
	assert.NoError(t, err)
	// check is id of checked by id is same with created user
	assert.Equal(t, user.ID, newUser.ID)
	assert.Equal(t, user.Name, newUser.Name)
	assert.Equal(t, user.Email, newUser.Email)
}

func TestUpdateUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_update_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create a user entity for testing
	newUser := &user_entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	userRepository.Create(newUser)

	newUserUpdate := &user_entity.User{
		ID:       1,
		Name:     "updated User",
		Email:    "updated@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	userRepository.Update(newUserUpdate)

	user, err := userRepository.ReadByID(newUser.ID)

	// Check if there is no error
	assert.NoError(t, err)
	// check is id of checked by id is same with created user
	assert.Equal(t, user.ID, newUserUpdate.ID)
	assert.Equal(t, user.Name, newUserUpdate.Name)
	assert.Equal(t, user.Email, newUserUpdate.Email)
}

func TestDeleteUser(t *testing.T) {
	// Setup database
	db, err := sqlite_configs.SetupDB("test_delete_user.db")
	if err != nil {
		t.Fatal(err)
	}

	// Defer closing the underlying database connection
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create a new instance of UserRepositoryStruct
	userRepository := user_repository_implementation.NewUserRepository(db)

	// Create a user entity for testing
	newUser := &user_entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
		Active:   true,
		Role:     1,
	}

	// Create the user
	createdUser, err := userRepository.Create(newUser)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	// Delete the user
	err = userRepository.SoftDelete(createdUser.ID)
	if err != nil {
		t.Fatalf("error deleting user: %v", err)
	}

	// Check if there is no error
	assert.NoError(t, err)

	// Attempt to read the deleted user by ID
	deletedUser, err := userRepository.ReadByID(createdUser.ID)

	// Check if the deleted user is nil
	assert.Nil(t, deletedUser)

	// Check if the error indicates that the user was not found
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

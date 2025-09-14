package test

import (
	"strings"
	"testing"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	repository_impl "github.com/celpung/gocleanarch/application/user/impl/repository"
	usecase_impl "github.com/celpung/gocleanarch/application/user/impl/usecase"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/*
===============================================================================
Test Execution Guide (Windows / macOS / Linux)

1) Install dependencies at the project root:
     go get github.com/glebarez/sqlite gorm.io/gorm github.com/stretchr/testify
     go mod tidy

2) Run all tests in this package from the folder containing this file:
     go test -v .

3) Run a specific test using a regex:
     go test -v -run ^TestUsecase_Create_ShouldHashPasswordAndPersist$ .
     go test -v -run ^TestUsecase_Update_NoChanges_ReturnsCurrentWithPasswordBlank$ .

4) Run the entire repository test suite from the project root:
     go test -v ./...

5) Coverage and race detection (optional):
     go test -race -cover .
     go test -coverprofile=cover.out . && go tool cover -html=cover.out

Notes:
- The User repository is backed by an in-memory SQLite database for isolation
  and speed. Schema is migrated for each test run.
- The PasswordService is used to hash and verify credentials. The tests avoid
  asserting the exact hash value and instead verify behavior via verification.
- Login success is not covered here to avoid coupling to JWT configuration.
  Error branches for wrong password and inactive user are validated.
===============================================================================
*/

/*
setupTestDB creates an isolated in-memory SQLite database and migrates the
schema required by the tests. Each test should call this to ensure a clean, independent state.
*/
// func setupTestDB(t *testing.T) *gorm.DB {
// 	t.Helper()

// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	require.NoError(t, err, "failed to open in-memory SQLite")

// 	/* Ensure the schema exists for all tests. The model should include DeletedAt so that soft deletes are correctly handled by GORM. */
// 	require.NoError(t, db.AutoMigrate(&model.User{}), "failed to auto-migrate schema")

// 	return db
// }

/*
newUsecase wires the use case with a real repository implementation and
concrete password and JWT services. The zero-value JWT service is sufficient
for tests that do not generate tokens.
*/
func newUsecase(t *testing.T) (*usecase_impl.UserUsecaseStruct, *gorm.DB) {
	t.Helper()

	db := setupTestDB(t)
	repo := repository_impl.NewUserRepository(db)

	/* PasswordService provides HashPassword and VerifyPassword. A zero-value instance is sufficient for tests since it is stateless. */
	ps := &auth.PasswordService{}

	/* JwtService is provided to satisfy the use case dependency. The tests here do not invoke token generation on success paths to avoid external config. */
	js := &auth.JwtService{}

	uc := &usecase_impl.UserUsecaseStruct{
		Repo:            repo,
		PasswordService: ps,
		JWTService:      js,
	}
	return uc, db
}

/*
helper user entity constructor for use case input. Password is plain here;
hashing is performed within the use case Create method.
*/
func makeEntityUser(name, email, plainPassword string, active bool, role uint) *entity.User {
	return &entity.User{
		Name:     name,
		Email:    email,
		Password: plainPassword,
		Active:   active,
		Role:     role,
	}
}

/* small pointer helpers for partial update payloads */
func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
func ptrUint(u uint) *uint       { return &u }

/*
TestUsecase_Create_ShouldHashPasswordAndPersist verifies that Create
hashes the incoming password and persists the user record.
*/
func TestUsecase_Create_ShouldHashPasswordAndPersist(t *testing.T) {
	uc, _ := newUsecase(t)

	in := makeEntityUser("Alice", "alice@ex.com", "secret123", true, 1)
	out, err := uc.Create(in)
	require.NoError(t, err, "create should not error")
	require.NotEmpty(t, out.ID, "expected ID to be set after create")

	/* Verify the stored password is a hash and validates against the original plaintext, using the repository path that exposes the password column. */
	stored, err := uc.Repo.ReadByEmailPrivate("alice@ex.com")
	require.NoError(t, err)
	require.NotEqual(t, "secret123", stored.Password, "stored password should not equal plaintext")

	err = uc.PasswordService.VerifyPassword(stored.Password, "secret123")
	require.NoError(t, err, "stored hash should verify with original plaintext")
}

/*
TestUsecase_Read_ReturnsEntitySlice verifies that Read returns entities
mapped from the repository projection and does not include sensitive fields.
*/
func TestUsecase_Read_ReturnsEntitySlice(t *testing.T) {
	uc, _ := newUsecase(t)

	_, err := uc.Create(makeEntityUser("Maria", "maria@ex.com", "pw", true, 1))
	require.NoError(t, err)
	_, err = uc.Create(makeEntityUser("Bob", "bob@ex.com", "pw", true, 1))
	require.NoError(t, err)

	list, err := uc.Read()
	require.NoError(t, err)
	require.Len(t, list, 2)

	/* Repository Read uses a public projection that omits the password column. The mapped entities should therefore have an empty Password field. */
	for _, e := range list {
		require.Empty(t, e.Password, "password should be empty in projected entity")
	}
}

/*
TestUsecase_ReadByID_ReturnsSingleEntity verifies that ReadByID maps the
repository model to the entity type correctly.
*/
func TestUsecase_ReadByID_ReturnsSingleEntity(t *testing.T) {
	uc, _ := newUsecase(t)

	created, err := uc.Create(makeEntityUser("Charlie", "charlie@ex.com", "pw", true, 2))
	require.NoError(t, err)

	got, err := uc.ReadByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)
	require.Equal(t, "Charlie", got.Name)
}

/*
TestUsecase_Update_NoChanges_ReturnsCurrentWithPasswordBlank exercises the
code path where no fields are provided for update. The use case should return
the current user state with the password sanitized.
*/
func TestUsecase_Update_NoChanges_ReturnsCurrentWithPasswordBlank(t *testing.T) {
	uc, _ := newUsecase(t)

	created, err := uc.Create(makeEntityUser("Diana", "diana@ex.com", "pw", true, 3))
	require.NoError(t, err)

	out, err := uc.Update(&entity.UpdateUserPayload{ID: created.ID})
	require.NoError(t, err)
	require.Equal(t, created.Name, out.Name)
	require.Equal(t, "", out.Password, "password should be blanked in the response")
}

/*
TestUsecase_Update_WriteZeroValues verifies that the map-based update path
writes zero values when requested through pointer fields in the payload.
*/
func TestUsecase_Update_WriteZeroValues(t *testing.T) {
	uc, _ := newUsecase(t)

	created, err := uc.Create(makeEntityUser("Eve", "eve@ex.com", "pw", true, 7))
	require.NoError(t, err)

	payload := &entity.UpdateUserPayload{
		ID:     created.ID,
		Name:   ptrString("Eve Zero"),
		Active: ptrBool(false), // request setting to false
		Role:   ptrUint(0),     // request setting to zero
	}
	out, err := uc.Update(payload)
	require.NoError(t, err)
	require.Equal(t, "Eve Zero", out.Name)

	/* Confirm persisted values through repository read that uses projection. */
	got, err := uc.Repo.ReadByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, "Eve Zero", got.Name)
	require.False(t, got.Active)
	require.EqualValues(t, 0, got.Role)
}

/*
TestUsecase_SoftDelete verifies that SoftDelete masks the record from subsequent default queries.
*/
func TestUsecase_SoftDelete(t *testing.T) {
	uc, _ := newUsecase(t)

	created, err := uc.Create(makeEntityUser("Frank", "frank@ex.com", "pw", true, 1))
	require.NoError(t, err)

	err = uc.SoftDelete(created.ID)
	require.NoError(t, err)

	_, err = uc.Repo.ReadByID(created.ID)
	require.Equal(t, gorm.ErrRecordNotFound, err, "expected soft-deleted record to be hidden")
}

/*
TestUsecase_Login_WrongPassword validates that providing an incorrect password
fails prior to any token generation and returns a descriptive error.
*/
func TestUsecase_Login_WrongPassword(t *testing.T) {
	uc, _ := newUsecase(t)

	_, err := uc.Create(makeEntityUser("Greg", "greg@ex.com", "right-pass", true, 1))
	require.NoError(t, err)

	token, err := uc.Login("greg@ex.com", "wrong-pass")
	require.Error(t, err)
	require.Empty(t, token)
	require.True(t, strings.Contains(err.Error(), "wrong password"))
}

/*
TestUsecase_Login_InactiveUser validates that login attempts are rejected for
inactive accounts prior to password verification or token generation.
*/
func TestUsecase_Login_InactiveUser(t *testing.T) {
	uc, _ := newUsecase(t)

	created, err := uc.Create(makeEntityUser("Hanna", "hanna@ex.com", "pw", true, 1))
	require.NoError(t, err)

	/* Mark the user inactive using the repository partial update. */
	_, err = uc.Repo.UpdateFields(created.ID, map[string]interface{}{"active": false})
	require.NoError(t, err)

	token, err := uc.Login("hanna@ex.com", "pw")
	require.Error(t, err)
	require.Empty(t, token)
	require.True(t, strings.Contains(err.Error(), "user not active"))
}

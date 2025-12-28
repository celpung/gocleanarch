package usecase_test

import (
	"context"
	"testing"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/celpung/gocleanarch/internal/infra/persistence"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeIDGen struct {
	id  string
	err error
}

func (f fakeIDGen) NewID() (string, error) { return f.id, f.err }

type fakeHasher struct {
	err error
}

func (f fakeHasher) Hash(plain string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return "hashed:" + plain, nil
}

func (f fakeHasher) Compare(_, _ string) error { return nil }

type fakeJWT struct{}

func (fakeJWT) Generate(_, _, _ string) (string, error) { return "token", nil }

var _ dependencies.JwtGenerator = (*fakeJWT)(nil)
var _ dependencies.PasswordHasher = (*fakeHasher)(nil)
var _ dependencies.IDGenerator = (*fakeIDGen)(nil)

func newTestUsecase(t *testing.T) (usecaseport.UserUsecase, repository.UserRepository, *gorm.DB) {
	t.Helper()

	db := setupTestDB(t)
	repo := persistence.NewUserRepository(db)

	uc := usecase.NewUserUsecase(
		repo,
		fakeIDGen{id: "id-123"},
		fakeHasher{},
		fakeJWT{},
		typograph.Typograph{},
	)

	return uc, repo, db
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	return db
}

func TestRegister_NormalizesAndHashes(t *testing.T) {
	uc, repo, _ := newTestUsecase(t)

	input := entity.User{
		Name:     " alice smith ",
		Email:    " Alice@example.com ",
		Role:     " user ",
		Password: "secret",
	}

	if err := uc.Register(context.Background(), input); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	got, err := repo.FindByID(context.Background(), "id-123")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	if got.Name != "Alice Smith" {
		t.Fatalf("expected normalized name, got %q", got.Name)
	}
	if got.Email != "Alice@example.com" {
		t.Fatalf("expected trimmed email, got %q", got.Email)
	}
	if got.Role != "user" {
		t.Fatalf("expected trimmed role, got %q", got.Role)
	}
	if got.Password != "hashed:secret" {
		t.Fatalf("expected hashed password, got %q", got.Password)
	}
}

func TestChangePassword_UsesAuthenticatedUser(t *testing.T) {
	uc, repo, _ := newTestUsecase(t)

	seed := entity.User{
		ID:       "id-1",
		Name:     "User",
		Email:    "user@example.com",
		Password: "hashed:old",
		Role:     "user",
	}
	if err := repo.Create(context.Background(), seed); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	if err := uc.ChangePassword(context.Background(), "id-1", "new-pass"); err != nil {
		t.Fatalf("ChangePassword returned error: %v", err)
	}

	updated, err := repo.FindByID(context.Background(), "id-1")
	if err != nil {
		t.Fatalf("FindByID after update: %v", err)
	}
	if updated.Password != "hashed:new-pass" {
		t.Fatalf("expected hashed new password, got %q", updated.Password)
	}
}

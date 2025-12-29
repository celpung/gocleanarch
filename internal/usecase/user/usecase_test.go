package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
	userdto "github.com/celpung/gocleanarch/internal/usecase/user/dto"
	"github.com/celpung/gocleanarch/internal/usecase/user/repository/memory"
	"github.com/celpung/gocleanarch/pkg/typograph"
)

type fakeIDGen struct {
	id  string
	err error
}

func (f fakeIDGen) NewID() (string, error) { return f.id, f.err }

type fakeHasher struct {
	hashErr    error
	compareErr error
	lastHash   string
}

func (f *fakeHasher) Hash(plain string) (string, error) {
	f.lastHash = plain
	if f.hashErr != nil {
		return "", f.hashErr
	}
	return "hashed:" + plain, nil
}

func (f fakeHasher) Compare(_, _ string) error {
	if f.compareErr != nil {
		return f.compareErr
	}
	return nil
}

type fakeJWT struct{}

func (fakeJWT) Generate(_, _, _ string) (string, error) { return "token", nil }

var _ dependencies.JwtGenerator = (*fakeJWT)(nil)
var _ dependencies.PasswordHasher = (*fakeHasher)(nil)
var _ dependencies.IDGenerator = (*fakeIDGen)(nil)

func newUserService(t *testing.T) (usecaseport.UserUsecase, *fakeHasher, repository.UserRepository) {
	t.Helper()

	repo := memory.NewUserRepository()
	hasher := &fakeHasher{}
	uc := usecase.NewUserUsecase(
		repo,
		fakeIDGen{id: "id-123"},
		hasher,
		fakeJWT{},
		typograph.Typograph{},
	)
	return uc, hasher, repo
}

func TestCreate_NormalizesAndHashes(t *testing.T) {
	uc, hasher, _ := newUserService(t)

	req := userdto.CreateUserRequest{
		Name:     " alice smith ",
		Email:    " Alice@example.com ",
		Role:     " user ",
		Password: "secret123",
	}

	resp, err := uc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if resp.ID != "id-123" {
		t.Fatalf("expected id %q, got %q", "id-123", resp.ID)
	}
	if resp.Name != "Alice Smith" {
		t.Fatalf("expected normalized name, got %q", resp.Name)
	}
	if resp.Email != "alice@example.com" {
		t.Fatalf("expected lowercased email, got %q", resp.Email)
	}
	if resp.Role != "user" {
		t.Fatalf("expected trimmed role, got %q", resp.Role)
	}
	if hasher.lastHash != "secret123" {
		t.Fatalf("expected password to be hashed, got %q", hasher.lastHash)
	}
}

func TestCreate_DuplicateEmail(t *testing.T) {
	uc, _, _ := newUserService(t)

	first := userdto.CreateUserRequest{
		Name:     "User",
		Email:    "user@example.com",
		Password: "secret123",
		Role:     "user",
	}
	if _, err := uc.Create(context.Background(), first); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err := uc.Create(context.Background(), first)
	if !errors.Is(err, entity.ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	uc, _, _ := newUserService(t)

	name := "Bob"
	_, err := uc.Update(context.Background(), userdto.UpdateUserRequest{
		ID:   "missing-id",
		Name: &name,
	})

	if !errors.Is(err, entity.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestChangePassword_Hashes(t *testing.T) {
	uc, hasher, repo := newUserService(t)

	// seed
	seed := userdto.CreateUserRequest{
		Name:     "User",
		Email:    "user@example.com",
		Password: "initialpass",
		Role:     "user",
	}
	created, err := uc.Create(context.Background(), seed)
	if err != nil {
		t.Fatalf("seed create failed: %v", err)
	}

	if err := uc.ChangePassword(context.Background(), created.ID, "new-pass"); err != nil {
		t.Fatalf("ChangePassword returned error: %v", err)
	}

	if hasher.lastHash != "new-pass" {
		t.Fatalf("expected ChangePassword to hash new password, got %q", hasher.lastHash)
	}

	u, err := repo.FindByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if u.Password != "hashed:new-pass" {
		t.Fatalf("expected stored hash, got %q", u.Password)
	}
}

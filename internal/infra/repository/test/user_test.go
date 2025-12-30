package test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/celpung/gocleanarch/internal/domain/entity"
	errs "github.com/celpung/gocleanarch/internal/domain/errors"
	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/celpung/gocleanarch/internal/infra/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	// Temp file per test (lebih stabil daripada :memory:)
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}

	return db, cleanup
}

func TestUserRepository_SQLite_Create_Success(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	err := repo.Create(context.Background(), entity.User{
		ID:       "11111111-2222-3333-4444-555555555555",
		Name:     "Edo",
		Email:    "edo@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestUserRepository_SQLite_Create_DuplicateEmail(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	err := repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "Edo",
		Email:    "dup@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	err = repo.Create(context.Background(), entity.User{
		ID:       "2",
		Name:     "Edo2",
		Email:    "dup@mail.com",
		Role:     "ADMIN",
		Password: "hashed2",
	})

	if err != errs.ErrEmailExists {
		t.Fatalf("expected ErrEmailExists, got: %v", err)
	}
}

func TestUserRepository_SQLite_Lists_Success_OrderByNameAsc(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{ID: "1", Name: "B", Email: "b@mail.com", Role: "ADMIN", Password: "x"})
	_ = repo.Create(context.Background(), entity.User{ID: "2", Name: "A", Email: "a@mail.com", Role: "ADMIN", Password: "y"})

	users, total, err := repo.Lists(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total=2, got: %d", total)
	}
	if len(users) != 2 {
		t.Fatalf("expected len=2, got: %d", len(users))
	}
	if users[0].Name != "A" {
		t.Fatalf("expected first user name A, got: %s", users[0].Name)
	}
}

func TestUserRepository_SQLite_FindByID_NotFound(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	u, err := repo.FindByID(context.Background(), "missing")
	if err != errs.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
	if u != nil {
		t.Fatalf("expected nil user, got: %+v", u)
	}
}

func TestUserRepository_SQLite_FindByEmail_NotFound(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	u, err := repo.FindByEmail(context.Background(), "missing@mail.com")
	if err != errs.ErrEmailNotFound {
		t.Fatalf("expected ErrEmailNotFound, got: %v", err)
	}
	if u != nil {
		t.Fatalf("expected nil user, got: %+v", u)
	}
}

func TestUserRepository_SQLite_FindByEmail_Success(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "Edo",
		Email:    "edo@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	u, err := repo.FindByEmail(context.Background(), "edo@mail.com")
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if u == nil || u.Email != "edo@mail.com" || u.Name != "Edo" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestUserRepository_SQLite_Update_InvalidInput_Nil(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	err := repo.Update(context.Background(), "1", nil)
	if err != errs.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestUserRepository_SQLite_Update_NoChanges_ReturnNil(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "Edo",
		Email:    "edo@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	// input kosong => updates map kosong => repo return nil
	err := repo.Update(context.Background(), "1", &entity.UpdateUser{})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestUserRepository_SQLite_Update_NotFound(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	name := "New Name"
	err := repo.Update(context.Background(), "missing", &entity.UpdateUser{
		Name: &name,
	})
	if err != errs.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_SQLite_Update_Success(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "Edo",
		Email:    "edo@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	newName := "Edo Updated"
	err := repo.Update(context.Background(), "1", &entity.UpdateUser{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	u, err := repo.FindByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
	if u.Name != "Edo Updated" {
		t.Fatalf("expected updated name, got: %s", u.Name)
	}
}

func TestUserRepository_SQLite_Update_DuplicateEmail(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "User1",
		Email:    "u1@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})
	_ = repo.Create(context.Background(), entity.User{
		ID:       "2",
		Name:     "User2",
		Email:    "u2@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	dup := "u2@mail.com"
	err := repo.Update(context.Background(), "1", &entity.UpdateUser{
		Email: &dup,
	})
	if err != errs.ErrEmailExists {
		t.Fatalf("expected ErrEmailExists, got: %v", err)
	}
}

func TestUserRepository_SQLite_Delete_NotFound(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	err := repo.Delete(context.Background(), "missing")
	if err != errs.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_SQLite_Delete_Success(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	_ = repo.Create(context.Background(), entity.User{
		ID:       "1",
		Name:     "Edo",
		Email:    "edo@mail.com",
		Role:     "ADMIN",
		Password: "hashed",
	})

	err := repo.Delete(context.Background(), "1")
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	_, err = repo.FindByID(context.Background(), "1")
	if err != errs.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound after delete, got: %v", err)
	}
}

package integration

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/celpung/gocleanarch/internal/infra/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := persistence.NewUserRepository(db)

	user := entity.User{
		ID:        "user-1",
		Name:      "Alice",
		Email:     "alice@example.com",
		Role:      "user",
		Password:  "hashed:secret",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	found, err := repo.FindByEmail(context.Background(), "alice@example.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}
	if found.ID != "user-1" || found.Email != "alice@example.com" {
		t.Fatalf("unexpected user found: %+v", found)
	}
	if found.CreatedAt.IsZero() || found.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be set")
	}
}

func TestUserRepository_DuplicateEmail(t *testing.T) {
	db := newTestDB(t)
	repo := persistence.NewUserRepository(db)

	user := entity.User{
		ID:       "user-1",
		Name:     "Alice",
		Email:    "alice@example.com",
		Role:     "user",
		Password: "hashed:secret",
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	if err := repo.Create(context.Background(), user); !errors.Is(err, entity.ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}

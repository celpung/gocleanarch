package migration

import (
	"fmt"

	"github.com/celpung/gocleanarch/internal/configs/db/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	// Tambah model lain di sini kalau sudah ada
	if err := db.AutoMigrate(
		&model.User{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	return nil
}

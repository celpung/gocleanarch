package migration

import (
	"fmt"

	"github.com/celpung/gocleanarch/internal_old/configs/db/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	if err := db.AutoMigrate(
		&model.Company{},
		&model.User{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	return nil
}

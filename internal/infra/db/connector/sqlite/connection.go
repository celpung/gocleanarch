package sqlite

import (
	"fmt"

	"github.com/celpung/gocleanarch/internal/infra/db/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupDB(dbname string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return db, nil
}

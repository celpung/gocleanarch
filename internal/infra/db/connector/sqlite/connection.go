package sqlite

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupDB(dbname string) (*gorm.DB, error) {
	if strings.TrimSpace(dbname) == "" {
		dbname = "app.db"
	}

	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	return db, nil
}

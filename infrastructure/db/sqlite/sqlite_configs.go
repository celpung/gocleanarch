package sqlite

import (
	"fmt"

	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SetupDB initializes a SQLite in-memory database for testing,
// creates the necessary table, and returns the DB connection.
func SetupDB(dbname string) (*gorm.DB, error) {
	// Initialize a SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// AutoMigrate creates the table based on the User struct
	if err := db.AutoMigrate(&user_entity.User{}); err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return db, nil
}

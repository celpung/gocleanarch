package test

import (
	"testing"

	"github.com/celpung/gocleanarch/infrastructure/db/model"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open in-memory SQLite database")

	/* Ensure the schema exists for all tests. The model should include DeletedAt so that soft deletes are correctly handled by GORM. */
	require.NoError(t, db.AutoMigrate(&model.User{}), "failed to auto-migrate schema")

	return db
}

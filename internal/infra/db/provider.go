package db

import (
	"fmt"
	"strings"

	"github.com/celpung/gocleanarch/internal/infra/db/connector/mysql"
	"github.com/celpung/gocleanarch/internal/infra/db/connector/sqlite"
	"github.com/celpung/gocleanarch/internal/infra/db/migration"
	"github.com/celpung/gocleanarch/internal/infra/environment"
	"gorm.io/gorm"
)

// Connector hides the concrete database connection details (dialect, driver, etc.).
type Connector interface {
	Connect(env environment.Environment) (*gorm.DB, error)
}

// Migrator abstracts schema migrations.
type Migrator interface {
	Migrate(db *gorm.DB) error
}

// MigratorFunc adapts a function to the Migrator interface.
type MigratorFunc func(db *gorm.DB) error

func (f MigratorFunc) Migrate(db *gorm.DB) error {
	return f(db)
}

// Provider chooses the correct connector by dialect and optionally runs migrations.
type Provider struct {
	connectors map[string]Connector
	migrator   Migrator
}

func DefaultProvider() Provider {
	return Provider{
		connectors: map[string]Connector{
			"mysql":  MySQLConnector{},
			"sqlite": SQLiteConnector{},
		},
		migrator: MigratorFunc(migration.Run),
	}
}

// Connect creates a DB connection using the configured dialect.
// If AUTO_MIGRATE is enabled, migrations run automatically.
func (p Provider) Connect(env environment.Environment) (*gorm.DB, error) {
	connector, ok := p.connectors[strings.ToLower(env.DB_DIALECT)]
	if !ok {
		return nil, fmt.Errorf("unsupported db dialect: %s", env.DB_DIALECT)
	}

	db, err := connector.Connect(env)
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %w", err)
	}

	// if env.AUTO_MIGRATE && p.migrator != nil {
	// 	if err := p.migrator.Migrate(db); err != nil {
	// 		return nil, fmt.Errorf("run migrations failed: %w", err)
	// 	}
	// }

	return db, nil
}

// MySQLConnector wires the existing MySQL adapter into the provider.
type MySQLConnector struct{}

func (MySQLConnector) Connect(env environment.Environment) (*gorm.DB, error) {
	cfg := mysql.Config{
		Username: env.DB_USERNAME,
		Password: env.DB_PASSWORD,
		Host:     env.DB_HOST,
		Port:     env.DB_PORT,
		Database: env.DB_NAME,
	}

	database, err := mysql.New(cfg)
	if err != nil {
		return nil, err
	}

	return database.DB, nil
}

// SQLiteConnector connects using the SQLite adapter.
type SQLiteConnector struct{}

func (SQLiteConnector) Connect(env environment.Environment) (*gorm.DB, error) {
	dbname := strings.TrimSpace(env.DB_NAME)
	if dbname == "" {
		dbname = "app.db"
	}

	db, err := sqlite.SetupDB(dbname)
	if err != nil {
		return nil, err
	}

	return db, nil
}

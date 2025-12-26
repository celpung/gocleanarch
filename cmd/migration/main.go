package main

import (
	"log"

	"github.com/celpung/gocleanarch/internal/configs/db/migration"
	"github.com/celpung/gocleanarch/internal/configs/db/mysql"
	"github.com/celpung/gocleanarch/internal/configs/environment"
)

func main() {
	env := environment.Load()

	cfg := mysql.Config{
		Username: env.DB_USERNAME,
		Password: env.DB_PASSWORD,
		Host:     env.DB_HOST,
		Port:     env.DB_PORT,
		Database: env.DB_NAME,
	}

	database, err := mysql.New(cfg)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	if err := migration.Migrate(database.DB); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration success ✅")
}

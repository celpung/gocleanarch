package main

import (
	"log"

	"github.com/celpung/gocleanarch/internal/infra/db"
	"github.com/celpung/gocleanarch/internal/infra/db/migration"
	"github.com/celpung/gocleanarch/internal/infra/environment"
)

func main() {
	env, err := environment.Load()
	if err != nil {
		log.Fatalf("env load failed: %v", err)
	}

	dbProvider := db.DefaultProvider()
	database, err := dbProvider.Connect(env)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	if err := migration.Migrate(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration success ✅")
}

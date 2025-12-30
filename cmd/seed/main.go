package main

import (
	"log"

	"github.com/celpung/gocleanarch/internal/infra/db"
	"github.com/celpung/gocleanarch/internal/infra/db/migration"
	"github.com/celpung/gocleanarch/internal/infra/db/seeder"
	"github.com/celpung/gocleanarch/internal/infra/environment"
)

func main() {
	env := environment.Load()

	dbProvider := db.DefaultProvider()
	database, err := dbProvider.Connect(env)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	if err := migration.Run(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := seeder.Seed(database); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("seed success ✅")
}

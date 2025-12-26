package main

import (
	"log"

	"github.com/celpung/gocleanarch/internal/configs/environment"
	"github.com/celpung/gocleanarch/internal/infra/db/mysql"
	"github.com/celpung/gocleanarch/internal/infra/db/seeder"
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

	if err := seeder.Seed(database.DB); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("seed success ✅")
}

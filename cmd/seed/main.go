package main

import (
	"log"

	"github.com/celpung/gocleanarch/app/infra/db/mysql"
	"github.com/celpung/gocleanarch/app/infra/db/seeder"
	"github.com/celpung/gocleanarch/app/infra/environment"
)

func main() {
	cfg := mysql.Config{
		Username: environment.Env.DB_USERNAME,
		Password: environment.Env.DB_PASSWORD,
		Host:     environment.Env.DB_HOST,
		Port:     environment.Env.DB_PORT,
		Database: environment.Env.DB_NAME,
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

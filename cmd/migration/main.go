package main

import (
	"log"

	"github.com/celpung/gocleanarch/app/infra/db/migration"
	"github.com/celpung/gocleanarch/app/infra/db/mysql"
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

	if err := migration.Migrate(database.DB); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration success ✅")
}

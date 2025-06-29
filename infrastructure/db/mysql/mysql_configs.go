package mysql

import (
	"database/sql"
	"fmt"

	slider_entity "github.com/celpung/gocleanarch/domain/slider/entity"
	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/celpung/gocleanarch/infrastructure/environment"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CreateDatabaseIfNotExists() error {
	dbUser := environment.Env.DB_USERNAME
	dbPassword := environment.Env.DB_PASSWORD
	dbHost := environment.Env.DB_HOST
	dbPort := environment.Env.DB_PORT
	dbName := environment.Env.DB_NAME

	// Connect to MySQL server without specifying a database
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort)
	sqlDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer sqlDB.Close()

	// Check if the database already exists
	var count int64
	if err := sqlDB.QueryRow("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&count); err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}

	// If the database doesn't exist, create it
	if count == 0 {
		if _, err := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

func ConnectDatabase() {
	dbUser := environment.Env.DB_USERNAME
	dbPassword := environment.Env.DB_PASSWORD
	dbHost := environment.Env.DB_HOST
	dbPort := environment.Env.DB_PORT
	dbName := environment.Env.DB_NAME

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}

func AutoMigrage() {
	ConnectDatabase()
	if migrateErr := DB.AutoMigrate(
		&user_entity.User{},
		&slider_entity.Slider{},
	); migrateErr != nil {
		panic(migrateErr)
	}
}

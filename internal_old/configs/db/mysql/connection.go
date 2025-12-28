package mysql

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New(cfg Config) (*Database, error) {
	if err := CreateDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	db, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func CreateDatabaseIfNotExists(cfg Config) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect mysql server failed: %w", err)
	}
	defer sqlDB.Close()

	var exists int
	query := `SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?`
	if err := sqlDB.QueryRow(query, cfg.Database).Scan(&exists); err != nil {
		return fmt.Errorf("check database existence failed: %w", err)
	}

	if exists == 0 {
		if _, err := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE `%s`", cfg.Database)); err != nil {
			return fmt.Errorf("create database failed: %w", err)
		}
	}

	return nil
}

func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %w", err)
	}

	return db, nil
}

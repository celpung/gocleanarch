package seeder

import (
	"fmt"

	"github.com/celpung/gocleanarch/app/infra/db/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	// contoh: seed 1 admin user (id tetap biar idempotent)
	const adminID = "00000000-0000-0000-0000-000000000001"
	const adminEmail = "admin@local.test"
	const adminName = "Admin"
	const adminPassword = "admin12345"

	// kalau sudah ada -> skip
	var count int64
	if err := db.Model(&model.User{}).Where("email = ?", adminEmail).Count(&count).Error; err != nil {
		return fmt.Errorf("check admin user failed: %w", err)
	}
	if count > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}

	admin := model.User{
		ID:       adminID,
		Name:     adminName,
		Email:    adminEmail,
		Password: string(hashed),
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("seed admin failed: %w", err)
	}

	return nil
}

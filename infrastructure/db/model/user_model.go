package model

import (
	"time"

	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string         `gorm:"unique"`
	Password  string         `gorm:"not null"`
	Active    bool           `gorm:"default:0"`
	Role      uint           `gorm:"not null;default:1"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func ToModel(e *user_entity.User) *User {
	var deletedAt gorm.DeletedAt
	if e.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *e.DeletedAt, Valid: true}
	} else {
		deletedAt = gorm.DeletedAt{Valid: false}
	}

	return &User{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Password:  e.Password,
		Active:    e.Active,
		Role:      e.Role,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ToEntity(m *User) *user_entity.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &user_entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		Active:    m.Active,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ToEntityList(models []*User) []*user_entity.User {
	entities := make([]*user_entity.User, 0, len(models))
	for _, m := range models {
		entities = append(entities, ToEntity(m))
	}
	return entities
}

func ToModelList(entities []*user_entity.User) []*User {
	models := make([]*User, 0, len(entities))
	for _, e := range entities {
		models = append(models, ToModel(e))
	}
	return models
}

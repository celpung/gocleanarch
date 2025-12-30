package port

import (
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

// dependencies
type Dependencies interface {
	NewID() (string, error)
	HashPassword(plain string) (string, error)
	ComparePassword(hashed string, plain string) error
	GenerateToken(id string, email string, role string) (string, error)
	ToTitleCase(s string) string
}

type (
	// user
	UserRepository = repository.UserRepository
	UserUsecase    = usecase.Usecase
)

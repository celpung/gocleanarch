package port

import (
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

// dependencies
type UUIDGenerator interface {
	NewID() (string, error)
}

type PasswordHasher interface {
	HashPassword(plain string) (string, error)
	ComparePassword(hashed string, plain string) error
}

type JWTGenerator interface {
	// GenerateJWT(userID, email, role string) (string, error)
	GenerateToken(id string, email string, role string) (string, error)
}

type Typograph interface {
	ToTitleCase(s string) string
}

type (
	// user
	UserRepository = repository.UserRepository
	UserUsecase    = usecase.Usecase
)

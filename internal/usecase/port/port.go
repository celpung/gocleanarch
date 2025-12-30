package port

import (
	"github.com/celpung/gocleanarch/internal/usecase/port/repository"
	"github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

type Dependencies interface {
	NewID() (string, error)
	HashPassword(plain string) (string, error)
	ComparePassword(hashed string, plain string) error
	GenerateJWT(userID, email, role string) (string, error)
}

type (
	// user
	UserRepository = repository.UserRepository
	UserUsecase    = usecase.Usecase
)

package user

import (
	"context"
	"time"

	domain "github.com/celpung/gocleanarch/internal_old/domain/user"
)

type Repository interface {
	Create(ctx context.Context, user domain.User) error
	Lists(ctx context.Context, offset, limit int) ([]domain.User, int64, error)
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id string, input *domain.UpdateUser) error
	Delete(ctx context.Context, id string) error
}

type IDGenerator interface {
	NewID() (string, error)
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashed string, plain string) error
}

type JWTGenerator interface {
	Generate(userID, email, role string) (string, error)
}

type AuthClaims interface {
	UserID() string
	Email() string
	Role() string
	IsExpired(leeway time.Duration) bool
	NotValidYet(leeway time.Duration) bool
}

type TokenVerifier interface {
	Verify(tokenStr string) (AuthClaims, error)
}

type Typograph interface {
	ToTitleCase(s string) string
}

type Usecase interface {
	Create(ctx context.Context, input CreateUserInput) (UserOutput, error)
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID, password string) error
	GetByID(ctx context.Context, id string) (UserOutput, error)
	List(ctx context.Context, page, limit int) ([]UserOutput, int64, error)
	Update(ctx context.Context, id string, updates UpdateUserInput) (UserOutput, error)
	Delete(ctx context.Context, id string) error
}

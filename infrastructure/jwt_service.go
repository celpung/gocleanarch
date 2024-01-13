package infrastructure

import (
	"os"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (js *JwtService) JWTGenerator(user entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
		"role":  user.Role,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

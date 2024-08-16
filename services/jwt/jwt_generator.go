package jwt_services

import (
	"os"

	"github.com/celpung/gocleanarch/entity"
	"github.com/golang-jwt/jwt"
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

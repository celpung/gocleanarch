package jwt_services

import (
	"time"

	"github.com/celpung/gocleanarch/configs/environment"
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
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(environment.Env.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

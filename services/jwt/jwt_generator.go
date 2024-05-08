package jwt_services

import (
	"os"

	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/golang-jwt/jwt"
)

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (js *JwtService) JWTGenerator(user user_entity.User) (string, error) {
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

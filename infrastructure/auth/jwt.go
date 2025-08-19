package auths

import (
	"time"

	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/golang-jwt/jwt/v4"
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
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(environment.Env.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package auth

import (
	"time"

	"github.com/celpung/gocleanarch/internal/configs/environment"
	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator struct{}

func (JwtGenerator) Generate(userID, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(environment.Load().JWT_TOKEN))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

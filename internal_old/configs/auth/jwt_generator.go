package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator struct {
	secret []byte
}

func NewJwtGenerator(secret string) JwtGenerator {
	return JwtGenerator{secret: []byte(secret)}
}

func (g JwtGenerator) Generate(userID, email, role string) (string, error) {
	if len(g.secret) == 0 {
		return "", errors.New("jwt secret is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(g.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

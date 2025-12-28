package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// JwtVerifier verifies JWT tokens using an HMAC secret.
type JwtVerifier struct {
	secret []byte
}

func NewJwtVerifier(secret string) JwtVerifier {
	return JwtVerifier{secret: []byte(secret)}
}

func (v JwtVerifier) Verify(tokenStr string, claims any) error {
	if len(v.secret) == 0 {
		return errors.New("jwt secret is empty")
	}

	jwtClaims, ok := claims.(jwt.Claims)
	if !ok {
		return fmt.Errorf("invalid claims type: %T", claims)
	}

	_, err := jwt.ParseWithClaims(tokenStr, jwtClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})
	return err
}

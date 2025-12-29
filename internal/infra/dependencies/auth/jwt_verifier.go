package auth

import (
	"errors"
	"time"

	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	"github.com/golang-jwt/jwt/v4"
)

// JwtVerifier verifies JWT tokens using an HMAC secret.
type JwtVerifier struct {
	secret []byte
}

func NewJwtVerifier(secret string) JwtVerifier {
	return JwtVerifier{secret: []byte(secret)}
}

func (v JwtVerifier) Verify(tokenStr string) (dependencies.AuthClaims, error) {
	if len(v.secret) == 0 {
		return nil, errors.New("jwt secret is empty")
	}

	var parsed jwtClaims
	_, err := jwt.ParseWithClaims(tokenStr, &parsed, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

type jwtClaims struct {
	ID       string `json:"id"`
	EmailVal string `json:"email"`
	RoleVal  string `json:"role"`
	jwt.RegisteredClaims
}

func (c *jwtClaims) UserID() string {
	return c.ID
}

func (c *jwtClaims) Email() string {
	return c.EmailVal
}

func (c *jwtClaims) Role() string {
	return c.RoleVal
}

func (c *jwtClaims) IsExpired(leeway time.Duration) bool {
	if c.ExpiresAt == nil {
		return false
	}
	return !c.ExpiresAt.After(time.Now().Add(-leeway))
}

func (c *jwtClaims) NotValidYet(leeway time.Duration) bool {
	if c.NotBefore == nil {
		return false
	}
	return c.NotBefore.After(time.Now().Add(leeway))
}

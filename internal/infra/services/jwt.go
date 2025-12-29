package services

import (
	"errors"
	"time"

	usecase "github.com/celpung/gocleanarch/internal/usecase/user"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) JWTService {
	return JWTService{secret: []byte(secret)}
}

func (s JWTService) Generate(userID, email, role string) (string, error) {
	if len(s.secret) == 0 {
		return "", errors.New("jwt secret is empty")
	}

	claims := jwtClaims{
		ID:       userID,
		EmailVal: email,
		RoleVal:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s JWTService) Verify(tokenStr string) (usecase.AuthClaims, error) {
	if len(s.secret) == 0 {
		return nil, errors.New("jwt secret is empty")
	}

	var parsed jwtClaims
	_, err := jwt.ParseWithClaims(tokenStr, &parsed, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
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

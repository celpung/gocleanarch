package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrMissingToken    = errors.New("auth: missing token")
	ErrInvalidToken    = errors.New("auth: invalid token")
	ErrInvalidClaims   = errors.New("auth: invalid claims")
	ErrInvalidAlgo     = errors.New("auth: invalid signing algorithm")
	ErrInvalidAuthHead = errors.New("auth: invalid authorization header")
)

type JwtClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`

	jwt.RegisteredClaims
}

type JwtService struct{}

// GenerateToken implements port.JWTGenerator.
func (js JwtService) GenerateToken(id string, email string, role string) (string, error) {
	panic("unimplemented")
}

func NewJwtService() *JwtService {
	return &JwtService{}
}
func (js *JwtService) GenerateToken(id string, email string, role string, secret string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	now := time.Now()

	claims := JwtClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

// ParseToken mem-parse JWT string menjadi JwtClaims dan memastikan token valid.
func (js *JwtService) ParseToken(tokenString string, secret string) (*JwtClaims, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrMissingToken
	}
	if strings.TrimSpace(secret) == "" {
		return nil, ErrInvalidToken // atau ErrMissingSecret
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, ErrInvalidAlgo
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || claims == nil {
		return nil, ErrInvalidClaims
	}

	// Validasi waktu exp/nbf/iat (lebih tegas)
	if err := claims.RegisteredClaims.Valid(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(claims.ID) == "" ||
		strings.TrimSpace(claims.Email) == "" ||
		strings.TrimSpace(claims.Role) == "" {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func (js *JwtService) ExtractBearerToken(authHeader string) (string, error) {
	h := strings.TrimSpace(authHeader)
	if h == "" {
		return "", ErrMissingToken
	}

	// Format: "Bearer <token>"
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 {
		return "", ErrInvalidAuthHead
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrInvalidAuthHead
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}

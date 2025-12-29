package dependencies

import "time"

// AuthClaims represents user identity extracted from a token.
type AuthClaims interface {
	UserID() string
	Email() string
	Role() string
	IsExpired(leeway time.Duration) bool
	NotValidYet(leeway time.Duration) bool
}

// TokenVerifier validates a token string and returns identity claims.
type TokenVerifier interface {
	Verify(tokenStr string) (AuthClaims, error)
}

package dependencies

// TokenVerifier validates a token string against provided claims.
// Using any keeps this decoupled from specific JWT claim implementations.
type TokenVerifier interface {
	Verify(tokenStr string, claims any) error
}

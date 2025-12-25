package port

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashed string, plain string) error
}

type IDGenerator interface {
	NewID() (string, error)
}

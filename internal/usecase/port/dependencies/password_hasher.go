package dependencies

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashed string, plain string) error
}

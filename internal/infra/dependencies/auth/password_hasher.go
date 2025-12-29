package auth

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func (BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (BcryptHasher) Compare(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

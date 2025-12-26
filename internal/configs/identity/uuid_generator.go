package identity

import "github.com/google/uuid"

type UUIDGenerator struct{}

func (UUIDGenerator) NewID() (string, error) {
	return uuid.NewString(), nil
}

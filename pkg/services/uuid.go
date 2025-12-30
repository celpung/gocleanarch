package services

import "github.com/google/uuid"

type UUIDService struct{}

func (UUIDService) NewID() (string, error) {
	return uuid.NewString(), nil
}

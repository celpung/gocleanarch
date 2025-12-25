package port

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrNoFieldsToUpdate = errors.New("no fields to update")

package port

import (
	"errors"
	"fmt"
)

// ValidationError represents a simple field-level validation error.
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

var (
	// user
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")

	//company
	ErrCompanyNameExists = errors.New("company name already exists")
	ErrCompanyNotFound   = errors.New("company not found")

	// global
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
	ErrNoCompanyAssociation = errors.New("company not found for provided company_id")
)

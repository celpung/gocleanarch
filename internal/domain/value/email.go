package value

import (
	"net/mail"
	"strings"

	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
)

// Email represents a validated, normalized email address.
type Email string

func NewEmail(raw string) (Email, error) {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	if normalized == "" {
		return "", apperrors.ErrEmailRequired
	}
	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", apperrors.ErrInvalidEmail
	}
	return Email(normalized), nil
}

func (e Email) String() string {
	return string(e)
}

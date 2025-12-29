package value

import (
	"strings"

	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
)

const MinPasswordLength = 8

// Password holds a validated, trimmed password.
type Password struct {
	value string
}

func NewPassword(raw string) (Password, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return Password{}, apperrors.ErrPasswordRequired
	}
	if len(normalized) < MinPasswordLength {
		return Password{}, apperrors.ErrWeakPassword
	}
	return Password{value: normalized}, nil
}

func (p Password) String() string {
	return p.value
}

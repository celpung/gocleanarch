package value

import (
	"strings"

	apperrors "github.com/celpung/gocleanarch/internal/domain/errors"
)

// Role captures the allowed user roles in the system.
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
	RoleSuper Role = "SUPER"
)

var allowedRoles = map[Role]struct{}{
	RoleUser:  {},
	RoleAdmin: {},
	RoleSuper: {},
}

func NewRole(raw string) (Role, error) {
	normalized := Role(strings.ToUpper(strings.TrimSpace(raw)))
	if normalized == "" {
		return "", apperrors.ErrRoleRequired
	}
	if _, ok := allowedRoles[normalized]; !ok {
		return "", apperrors.ErrInvalidRole
	}
	return normalized, nil
}

func (r Role) String() string {
	return string(r)
}

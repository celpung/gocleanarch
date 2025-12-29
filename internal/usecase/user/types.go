package user

import "time"

// CreateUserInput represents user creation input at the usecase boundary.
type CreateUserInput struct {
	Email    string
	Name     string
	Password string
	Role     string
}

// UpdateUserInput represents partial updates for a user.
type UpdateUserInput struct {
	Email    *string
	Name     *string
	Password *string
	Role     *string
}

// UserOutput is the outward-facing representation of a user from the use case layer.
// It excludes internal-only fields such as the password hash.
type UserOutput struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

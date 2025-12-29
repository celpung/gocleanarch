package dto

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

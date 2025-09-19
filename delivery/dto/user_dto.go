package dto

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	Role     string `json:"role" binding:"required" validate:"required"`
}

type UserUpdateRequest struct {
	ID       string  `json:"id" binding:"required" validate:"required,uuid4"`
	Name     *string `json:"name" binding:"omitempty" validate:"omitempty"`
	Email    *string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=8" validate:"omitempty,min=8"`
	Active   *bool   `json:"active" binding:"omitempty" validate:"omitempty"`
	Role     *string `json:"role" binding:"omitempty" validate:"omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
	Role   string `json:"role"`
}

package dto

type CreateUserRequest struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	CompanyID *string `json:"company_id,omitempty"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	CompanyID *string `json:"company_id"`
	IsActive  *bool   `json:"is_active"`
}

type UserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CompanyID *string `json:"company_id,omitempty"`
	IsActive  bool    `json:"is_active"`
}

type ListUsersResponse struct {
	Data []UserResponse `json:"data"`
	Meta PagingMeta     `json:"meta"`
}

type PagingMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
}

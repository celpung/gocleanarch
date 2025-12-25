package dto

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListUsersResponse struct {
	Data []UserItem `json:"data"`
	Meta PagingMeta `json:"meta"`
}

type UserItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PagingMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

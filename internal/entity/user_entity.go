package entity

type User struct {
	ID       string
	Name     string
	Email    string
	Role     string
	Password string
}

type UpdateUser struct {
	Name     *string
	Email    *string
	Role     *string
	Password *string
}

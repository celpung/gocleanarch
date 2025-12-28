package entity

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      string
	CompanyID string
	IsActive  bool
}

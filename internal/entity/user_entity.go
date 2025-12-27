package entity

// User represents an application user stored in the database.
type User struct {
	ID        string `gorm:"type:char(36);primaryKey;column:id"`
	Name      string `gorm:"size:120;not null;column:name"`
	Email     string `gorm:"size:120;not null;uniqueIndex;column:email"`
	Password  string `gorm:"size:255;not null;column:password"`
	Role      string ``
	CompanyID string `gorm:"type:char(36);column:company_id"`
	IsActive  bool   `gorm:"not null;default:true;column:is_active"`
}

func (User) TableName() string {
	return "users"
}

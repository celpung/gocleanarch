package model

// Company is the persistence representation of companies table.
// It should only be used by repository adapters.
type Company struct {
	ID       string `gorm:"type:char(36);primaryKey;column:id"`
	Name     string `gorm:"size:120;not null;uniqueIndex;column:name"`
	Address  string `gorm:"size:255;column:address"`
	Phone    string `gorm:"size:50;column:phone"`
	IsActive bool   `gorm:"not null;default:true;column:is_active"`
}

func (Company) TableName() string {
	return "companies"
}

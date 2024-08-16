package user_repository_implementation

import (
	user_repository "github.com/celpung/gocleanarch/domain/user/repository"
	"github.com/celpung/gocleanarch/entity"
	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	DB *gorm.DB
}

// Create implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Create(user *entity.User) (*entity.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Delete implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Delete(userID uint) error {
	if err := r.DB.Delete(&entity.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

// Read implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Read() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ReadByID implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByID(userID uint) (*entity.User, error) {
	var user entity.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ReadByEmail implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Update(user *entity.User) (*entity.User, error) {
	if err := r.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserRepositry(db *gorm.DB) user_repository.UserRepositoryInterface {
	return &UserRepositoryStruct{DB: db}
}

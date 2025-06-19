package user_repository_implementation

import (
	user_entity "github.com/celpung/gocleanarch/domain/user/entity"
	user_repository "github.com/celpung/gocleanarch/domain/user/repository"
	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	DB *gorm.DB
}

// Create implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Create(user *user_entity.User) (*user_entity.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Read implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Read() ([]*user_entity.User, error) {
	var users []*user_entity.User
	if err := r.selectUserData(r.DB).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ReadByID implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByID(userID uint) (*user_entity.User, error) {
	var user user_entity.User
	if err := r.selectUserData(r.DB).First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ReadByEmail implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByEmail(email string, isLogin bool) (*user_entity.User, error) {
	var user user_entity.User
	if isLogin {
		if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.selectUserData(r.DB).Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// Update implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Update(user *user_entity.User) (*user_entity.User, error) {
	if err := r.DB.Model(&user_entity.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete implements user_repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Delete(userID uint) error {
	if err := r.DB.Delete(&user_entity.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryStruct) selectUserData(db *gorm.DB) *gorm.DB {
	return db.Select("ID, Name, Email, Active, Role")
}

func NewUserRepositry(db *gorm.DB) user_repository.UserRepositoryInterface {
	return &UserRepositoryStruct{DB: db}
}

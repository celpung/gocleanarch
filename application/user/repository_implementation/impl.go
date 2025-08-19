package repository_implementation

import (
	"github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/celpung/gocleanarch/domain/user/repository"
	"github.com/celpung/gocleanarch/infrastructure/db/model"
	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	DB *gorm.DB
}

// Create implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Create(user *entity.User) (*entity.User, error) {
	usr := model.ToModel(user)
	if err := r.DB.Create(usr).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(usr), nil
}

// Read implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Read() ([]*entity.User, error) {
	var users []*model.User
	if err := r.selectUserData(r.DB).Find(&users).Error; err != nil {
		return nil, err
	}

	return model.ToEntityList(users), nil
}

// ReadByID implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByID(userID uint) (*entity.User, error) {
	var user model.User
	if err := r.selectUserData(r.DB).First(&user, userID).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(&user), nil
}

// ReadByEmail implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) ReadByEmail(email string, isLogin bool) (*entity.User, error) {
	var user model.User
	if isLogin {
		if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.selectUserData(r.DB).Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	}

	return model.ToEntity(&user), nil
}

// Update implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) Update(user *entity.User) (*entity.User, error) {
	User := model.ToModel(user)
	if err := r.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(User).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(User), nil
}

// Delete implements repository.UserRepositoryInterface.
func (r *UserRepositoryStruct) SoftDelete(userID uint) error {
	if err := r.DB.Delete(&model.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryStruct) selectUserData(db *gorm.DB) *gorm.DB {
	return db.Select("ID, Name, Email, Active, Role")
}

func NewUserRepository(db *gorm.DB) repository.UserRepositoryInterface {
	return &UserRepositoryStruct{DB: db}
}

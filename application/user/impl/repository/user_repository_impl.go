package repository_impl

import (
	"github.com/celpung/gocleanarch/application/user/domain/repository"
	"github.com/celpung/gocleanarch/infrastructure/db/model"
	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	DB *gorm.DB
}

func (r *UserRepositoryStruct) Create(m *model.User) (*model.User, error) {
	if err := r.DB.Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *UserRepositoryStruct) Read(page, limit uint) ([]*model.User, int64, error) {
	var (
		users []*model.User
		total int64
	)

	const defaultLimit uint = 10
	const maxLimit uint = 100

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := int((page - 1) * limit)

	base := r.DB.Model(&model.User{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.selectUserData(base.Session(&gorm.Session{})).
		Order("created_at DESC").
		Offset(offset).
		Limit(int(limit)).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepositoryStruct) ReadByID(userID string) (*model.User, error) {
	user := &model.User{}
	if err := r.selectUserData(r.DB).
		First(user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryStruct) ReadByEmailPublic(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.selectUserData(r.DB).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryStruct) ReadByEmailPrivate(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.DB.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryStruct) Update(user *model.User) (*model.User, error) {
	if err := r.DB.Model(model.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryStruct) UpdateFields(id string, fields map[string]interface{}) (*model.User, error) {
	tx := r.DB.Model(&model.User{}).Where("id = ?", id).Updates(fields)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var m model.User
	if err := r.DB.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserRepositoryStruct) SoftDelete(userID string) error {
	if err := r.DB.
		Where("id = ?", userID). // ← pakai "id = ?"
		Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryStruct) selectUserData(db *gorm.DB) *gorm.DB {
	return db.Select("ID, Name, Email, Active, Role")
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryStruct{DB: db}
}

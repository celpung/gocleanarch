package usecase_implementation

import (
	"errors"

	"github.com/celpung/gocleanarch/domain/user/entity"
	"github.com/celpung/gocleanarch/domain/user/repository"
	"github.com/celpung/gocleanarch/domain/user/usecase"
	"github.com/celpung/gocleanarch/infrastructure/auths"
)

type UserUsecaseStruct struct {
	UserRepository  repository.UserRepositoryInterface
	PasswordService *auths.PasswordService
	JWTService      *auths.JwtService
}

// Create implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Create(user *entity.User) (*entity.User, error) {
	// hashing password
	hashedPassword, err := u.PasswordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// set the hashed password into new user password
	user.Password = hashedPassword

	// perform create user
	user, userErr := u.UserRepository.Create(user)
	if userErr != nil {
		return nil, userErr
	}

	return user, nil
}

// Read implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Read() ([]*entity.User, error) {
	return u.UserRepository.Read()
}

// ReadByID implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) ReadByID(userID uint) (*entity.User, error) {
	return u.UserRepository.ReadByID(userID)
}

// Update implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Update(user *entity.User) (*entity.User, error) {
	existingUser, err := u.UserRepository.ReadByID(user.ID)
	if err != nil {
		return nil, err
	}

	// Update only the non-zero fields
	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		hashedPassword, err := u.PasswordService.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		existingUser.Password = hashedPassword
	}
	if user.Active {
		existingUser.Active = user.Active
	}
	if user.Role > 0 {
		existingUser.Role = user.Role
	}

	// Perform the update operation
	updatedUser, err := u.UserRepository.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// Delete implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) SoftDelete(userID uint) error {
	// perform soft delete user
	return u.UserRepository.SoftDelete(userID)
}

// Login implements usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Login(email, password string) (string, error) {
	// perform read user by email
	user, err := u.UserRepository.ReadByEmail(email, true)
	if err != nil {
		return "", err
	}

	// check is user active
	if !user.Active {
		return "", errors.New("user not active")
	}

	// verify hash password match plain password
	if err := u.PasswordService.VerifyPassword(user.Password, password); err != nil {
		return "", errors.New("wrong password")
	}

	// generate jwt token
	token, err := u.JWTService.JWTGenerator(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserUsecase(repository repository.UserRepositoryInterface, passwordServive *auths.PasswordService, jwtService *auths.JwtService) usecase.UserUsecaseInterface {
	return &UserUsecaseStruct{
		UserRepository:  repository,
		PasswordService: passwordServive,
		JWTService:      jwtService,
	}
}

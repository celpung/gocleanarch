package user_usecase_implementation

import (
	"errors"

	user_repository "github.com/celpung/gocleanarch/domain/user/repository"
	user_usecase "github.com/celpung/gocleanarch/domain/user/usecase"
	"github.com/celpung/gocleanarch/entity"
	jwt_services "github.com/celpung/gocleanarch/services/jwt"
	password_services "github.com/celpung/gocleanarch/services/password"
)

type UserUsecaseStruct struct {
	UserRepository  user_repository.UserRepositoryInterface
	PasswordService *password_services.PasswordService
	JWTService      *jwt_services.JwtService
}

// Create implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Create(user *entity.User) (*entity.User, error) {
	// hashing password
	hashedPassword, err := u.PasswordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// set the hashed password into new user password
	user.Password = hashedPassword

	// perform create user
	return u.UserRepository.Create(user)
}

// Delete implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Delete(userID uint) error {
	// perform delete user
	return u.UserRepository.Delete(userID)
}

// Read implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Read() ([]*entity.User, error) {
	// perform read all user
	return u.UserRepository.Read()
}

// ReadByID implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) ReadByID(userID uint) (*entity.User, error) {
	// perform read user by id
	return u.UserRepository.ReadByID(userID)
}

// Update implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Update(user *entity.User) (*entity.User, error) {
	// get user data by id
	userData, err := u.UserRepository.ReadByID(user.ID)
	if err != nil {
		return nil, err
	}

	// update the current user data into new user data as needed
	userData.Name = user.Name
	userData.Email = user.Email

	// perform update user
	return u.UserRepository.Update(userData)
}

// Login implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Login(email, password string) (string, error) {
	// perform read user by email
	user, err := u.UserRepository.ReadByEmail(email)
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

func NewUserUsecase(repository user_repository.UserRepositoryInterface, passwordServive *password_services.PasswordService, jwtService *jwt_services.JwtService) user_usecase.UserUsecaseInterface {
	return &UserUsecaseStruct{
		UserRepository:  repository,
		PasswordService: passwordServive,
		JWTService:      jwtService,
	}
}

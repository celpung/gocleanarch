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
func (u *UserUsecaseStruct) Create(user *entity.User) (*entity.UserHttpResponse, error) {
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

	userResponse := &entity.UserHttpResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Active: user.Active,
		Role:   user.Role,
	}

	return userResponse, nil
}

// Delete implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Delete(userID uint) error {
	// perform delete user
	return u.UserRepository.Delete(userID)
}

// Read implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Read() ([]*entity.UserHttpResponse, error) {
	// perform read all user
	user, err := u.UserRepository.Read()
	if err != nil {
		return nil, err
	}

	var userResponse []*entity.UserHttpResponse
	for _, v := range user {
		userResponse = append(userResponse, &entity.UserHttpResponse{
			ID:     v.ID,
			Name:   v.Name,
			Email:  v.Email,
			Active: v.Active,
			Role:   v.Role,
		})
	}
	return userResponse, nil

	// return u.UserRepository.Read()
}

// ReadByID implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) ReadByID(userID uint) (*entity.UserHttpResponse, error) {
	// perform read user by id
	user, userErr := u.UserRepository.ReadByID(userID)
	if userErr != nil {
		return nil, userErr
	}

	userResponse := &entity.UserHttpResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Active: user.Active,
		Role:   user.Role,
	}

	return userResponse, nil
}

// Update implements user_usecase.UserUsecaseInterface.
func (u *UserUsecaseStruct) Update(user *entity.UserUpdate) (*entity.UserHttpResponse, error) {
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

	userResponse := &entity.UserHttpResponse{
		ID:     updatedUser.ID,
		Name:   updatedUser.Name,
		Email:  updatedUser.Email,
		Active: updatedUser.Active,
		Role:   updatedUser.Role,
	}

	return userResponse, nil
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

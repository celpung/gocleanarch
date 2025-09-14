package delivery_impl

import (
	"net/http"

	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/delivery/dto"
	delivery "github.com/celpung/gocleanarch/delivery/gin/user"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
	"github.com/celpung/gocleanarch/infrastructure/validation"
	"github.com/gin-gonic/gin"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecase
}

func (d *UserDeliveryStruct) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data", "error": err.Error()})
		return
	}
	if err := validation.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	var e entity.User
	if err := mapper.CopyTo(&req, &e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map request", "error": err.Error()})
		return
	}

	user, err := d.UserUsecase.Create(&e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	var res dto.UserResponse
	if err := mapper.CopyTo(user, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map response", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Register success", "user": res})
}

func (d *UserDeliveryStruct) GetAllUserData(c *gin.Context) {
	users, err := d.UserUsecase.Read()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user data", "error": err.Error()})
		return
	}

	resp, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map response list", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success fetch user data", "users": resp})
}

func (d *UserDeliveryStruct) UpdateUser(c *gin.Context) {
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid update data", "error": err.Error()})
		return
	}
	if err := validation.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	var payload entity.UpdateUserPayload
	if err := mapper.CopyTo(&req, &payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map update payload", "error": err.Error()})
		return
	}

	user, err := d.UserUsecase.Update(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user", "error": err.Error()})
		return
	}

	var resp dto.UserResponse
	if err := mapper.CopyTo(user, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to map response", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": resp})
}

func (d *UserDeliveryStruct) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")

	if err := d.UserUsecase.SoftDelete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (d *UserDeliveryStruct) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid login data", "error": err.Error()})
		return
	}
	if err := validation.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "error": err.Error()})
		return
	}

	token, err := d.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Login failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login success", "token": token})
}

func NewUserDelivery(usecase usecase.UserUsecase) delivery.UserDelivery {
	return &UserDeliveryStruct{UserUsecase: usecase}
}

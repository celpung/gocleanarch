package user_delivery_implementation

import (
	"net/http"
	"strconv"

	"github.com/celpung/gocleanarch/delivery/dto"
	"github.com/celpung/gocleanarch/delivery/gin/user_delivery"
	"github.com/celpung/gocleanarch/domain/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecaseInterface
}

func (d *UserDeliveryStruct) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	userEntity := dto.UserCreateRequestDTO(&req)

	user, err := d.UserUsecase.Create(userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register success!",
		"user":    dto.UserResponseDTO(user),
	})
}

// GetAllUserData implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) GetAllUserData(c *gin.Context) {
	user, err := d.UserUsecase.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to bind login data!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success fetch user data!",
		"user":    dto.UserResponseListDTO(user),
	})
}

func (d *UserDeliveryStruct) UpdateUser(c *gin.Context) {
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to bind update data!",
			"error":   err.Error(),
		})
		return
	}

	userEntity := dto.UserUpdateRequestDTO(&req)

	user, err := d.UserUsecase.Update(userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update data!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully!",
		"user":    dto.UserResponseDTO(user),
	})
}

// DeleteUser implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
		return
	}

	if err := d.UserUsecase.SoftDelete(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// Login implements user_delivery.UserDeliveryInterface.
func (d *UserDeliveryStruct) Login(c *gin.Context) {
	var login dto.UserLoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to bind login data!",
			"error":   err.Error(),
		})
		return
	}

	// perform login to get token data
	token, err := d.UserUsecase.Login(login.Email, login.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Login failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   token,
	})
}

func NewUserDelivery(usecase usecase.UserUsecaseInterface) user_delivery.UserDeliveryInterface {
	return &UserDeliveryStruct{
		UserUsecase: usecase,
	}
}

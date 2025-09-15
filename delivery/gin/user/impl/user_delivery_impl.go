package delivery_impl

import (
	"net/http"
	"strconv"

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
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	// Ambil & normalisasi query
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	// Call usecase
	users, total, err := d.UserUsecase.Read(uint(page), uint(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user data",
			"error":   err.Error(),
		})
		return
	}

	// Map ke DTO
	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
		return
	}

	// Hitung total_page (ceil div)
	totalPage := (total + int64(limit) - 1) / int64(limit)

	// JSON shape sama seperti std/http kamu
	c.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"data": gin.H{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   totalPage,
		},
	})
}

func (d *UserDeliveryStruct) SearchUser(c *gin.Context) {
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := defaultPage
	limit := defaultLimit

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err != nil || v < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page parameter"})
			return
		} else {
			page = v
		}
	}
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err != nil || v < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit parameter"})
			return
		} else {
			limit = v
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	keyword := c.Query("q")

	users, total, err := d.UserUsecase.Search(uint(page), uint(limit), keyword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data not found"})
		return
	}

	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
		return
	}

	totalPage := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"data": gin.H{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   totalPage,
		},
	})
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

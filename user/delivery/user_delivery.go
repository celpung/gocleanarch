package delivery

import (
	"net/http"

	"github.com/celpung/gocleanarch/domain"
	"github.com/celpung/gocleanarch/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserDelivery struct {
	UserUsecase usecase.UserUsecase
}

func NewUserDelivery(userUseCase usecase.UserUsecase) *UserDelivery {
	return &UserDelivery{
		UserUsecase: userUseCase,
	}
}

func (ud *UserDelivery) Create(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Register failed",
		})
		return
	}

	if err := ud.UserUsecase.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register success",
	})
}

// Handler for retrieving users.
func (ud *UserDelivery) Read(ctx *gin.Context) {
	users, err := ud.UserUsecase.Read(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

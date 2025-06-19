package user_delivery

import "github.com/gin-gonic/gin"

type UserDeliveryInterface interface {
	Register(c *gin.Context)
	GetAllUserData(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Login(c *gin.Context)
}

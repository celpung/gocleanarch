package delivery

import "github.com/gin-gonic/gin"

type UserDelivery interface {
	Register(c *gin.Context)
	GetAllUserData(c *gin.Context)
	SearchUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Login(c *gin.Context)
}

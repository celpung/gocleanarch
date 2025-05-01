package crud_delivery

import "github.com/gin-gonic/gin"

type DeliveryInterface[T any] interface {
	Create(c *gin.Context)
	Read(c *gin.Context)
	ReadByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Search(c *gin.Context)
	Count(c *gin.Context)
}

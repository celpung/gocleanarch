package user_router

import (
	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	"github.com/celpung/gocleanarch/configs/role"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/gin/user_delivery/implementation"
	"github.com/celpung/gocleanarch/delivery/gin/user_delivery/middlewares"
	user_repository_implementation "github.com/celpung/gocleanarch/domain/user/repository/implementation"
	user_usecase_implementation "github.com/celpung/gocleanarch/domain/user/usecase/implementation"
	"github.com/celpung/gocleanarch/services"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	passwordService := services.NewPasswordService()
	jwtService := services.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql_configs.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	routes := r.Group("/users")
	{
		routes.POST("/register", delivery.Register)
		routes.POST("/login", delivery.Login)
		routes.GET("", middlewares.AuthMiddleware(role.Admin), delivery.GetAllUserData)
		routes.PATCH("", middlewares.AuthMiddleware(role.User), delivery.UpdateUser)
	}
}

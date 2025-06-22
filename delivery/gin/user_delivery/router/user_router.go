package user_router

import (
	user_repository_implementation "github.com/celpung/gocleanarch/application/user/repository"
	user_usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/gin/user_delivery/implementation"
	"github.com/celpung/gocleanarch/delivery/gin/user_delivery/middlewares"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql.DB)
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

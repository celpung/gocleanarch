package user_router

import (
	repository_implementation "github.com/celpung/gocleanarch/application/user/repository_implementation"
	usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase_implementation"
	delivery_implementation "github.com/celpung/gocleanarch/delivery/gin/user/implementation"
	"github.com/celpung/gocleanarch/delivery/gin/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()

	repository := repository_implementation.NewUserRepository(mysql.DB)
	usecase := usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_implementation.NewUserDelivery(usecase)

	routes := r.Group("/users")
	{
		routes.POST("/register", delivery.Register)
		routes.POST("/login", delivery.Login)
		routes.GET("", middleware.AuthMiddleware(role.Admin), delivery.GetAllUserData)
		routes.PATCH("", middleware.AuthMiddleware(role.User), delivery.UpdateUser)
	}
}

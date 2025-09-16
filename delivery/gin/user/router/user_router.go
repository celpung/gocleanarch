package user_router

import (
	repository_impl "github.com/celpung/gocleanarch/application/user/impl/repository"
	usecase_impl "github.com/celpung/gocleanarch/application/user/impl/usecase"
	delivery_impl "github.com/celpung/gocleanarch/delivery/gin/user/impl"
	"github.com/celpung/gocleanarch/delivery/gin/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()

	repository := repository_impl.NewUserRepository(mysql.DB)
	usecase := usecase_impl.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_impl.NewUserDelivery(usecase)

	routes := r.Group("/users")
	{
		routes.POST("/register", delivery.Register)
		routes.POST("/login", delivery.Login)
		routes.GET("", middleware.AuthMiddleware(middleware.Admin, middleware.Super), delivery.GetAllUserData)
		routes.GET("/search", middleware.AuthMiddleware(middleware.Admin), delivery.SearchUser)
		routes.PATCH("", middleware.AuthMiddleware(middleware.User), delivery.UpdateUser)
	}
}

package routers

import (
	"github.com/celpung/gocleanarch/configs"
	"github.com/celpung/gocleanarch/infrastructure"
	"github.com/celpung/gocleanarch/internal/handler"
	repository "github.com/celpung/gocleanarch/internal/repository/user"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	passwordSrv := infrastructure.NewPasswordService()
	jwtSrv := infrastructure.NewJwtService()

	// User routes
	userRepo := repository.NewUserRepository(configs.DB)
	userUseCase := usecase.NewUserUseCase(*userRepo, passwordSrv, jwtSrv)
	userHandler := handler.NewUserHandler(*userUseCase)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", middlewares.JWTMiddleware(configs.Admin), userHandler.Create)
		userRoutes.GET("", middlewares.JWTMiddleware(configs.Admin), userHandler.Read)
	}
}

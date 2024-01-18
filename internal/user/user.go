package user

import (
	"github.com/celpung/gocleanarch/configs"
	"github.com/celpung/gocleanarch/services"
	"github.com/celpung/gocleanarch/internal/user/delivery"
	"github.com/celpung/gocleanarch/internal/user/repository"
	"github.com/celpung/gocleanarch/internal/user/usecase"
	"github.com/celpung/gocleanarch/middlewares"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	passwordSrv := services.NewPasswordService()
	jwtSrv := services.NewJwtService()

	userRepo := repository.NewUserRepository(configs.DB)
	userUseCase := usecase.NewUserUseCase(userRepo, passwordSrv, jwtSrv)
	userDelivery := delivery.NewUserDelivery(*userUseCase)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", userDelivery.Create)
		userRoutes.GET("", middlewares.JWTMiddleware(configs.Admin), userDelivery.Read)
	}
}

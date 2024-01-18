package user

import (
	"github.com/celpung/gocleanarch/configs"
	"github.com/celpung/gocleanarch/infrastructure"
	"github.com/celpung/gocleanarch/middlewares"
	"github.com/celpung/gocleanarch/user/delivery"
	"github.com/celpung/gocleanarch/user/repository"
	"github.com/celpung/gocleanarch/user/usecase"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	passwordSrv := infrastructure.NewPasswordService()
	jwtSrv := infrastructure.NewJwtService()

	userRepo := repository.NewUserRepository(configs.DB)
	userUseCase := usecase.NewUserUseCase(userRepo, passwordSrv, jwtSrv)
	userDelivery := delivery.NewUserDelivery(*userUseCase)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", userDelivery.Create)
		userRoutes.GET("", middlewares.JWTMiddleware(configs.Admin), userDelivery.Read)
	}
}

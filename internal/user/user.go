package user

import (
	"github.com/celpung/gocleanarch/configs"
	"github.com/celpung/gocleanarch/internal/user/delivery"
	repositoryimplementation "github.com/celpung/gocleanarch/internal/user/repository/implementation"
	usecaseimplementation "github.com/celpung/gocleanarch/internal/user/usecase/implementation"
	"github.com/celpung/gocleanarch/middlewares"
	"github.com/celpung/gocleanarch/services"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	passwordSrv := services.NewPasswordService()
	jwtSrv := services.NewJwtService()

	userRepo := repositoryimplementation.NewUserRepository(configs.DB)
	userUseCase := usecaseimplementation.NewUserUseCase(userRepo, passwordSrv, jwtSrv)
	userDelivery := delivery.NewUserDelivery(userUseCase)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", userDelivery.Create)
		userRoutes.GET("", middlewares.JWTMiddleware(configs.Admin), userDelivery.Read)
	}
}

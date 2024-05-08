package user_router

import (
	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	middlewares "github.com/celpung/gocleanarch/configs/middlewares/gin"
	"github.com/celpung/gocleanarch/configs/role"
	user_delivery_implementation "github.com/celpung/gocleanarch/domain/user/delivery/gin/implementation"
	user_repository_implementation "github.com/celpung/gocleanarch/domain/user/repository/implementation"
	user_usecase_implementation "github.com/celpung/gocleanarch/domain/user/usecase/implementation"
	jwt_services "github.com/celpung/gocleanarch/services/jwt"
	password_services "github.com/celpung/gocleanarch/services/password"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	passwordService := password_services.NewPasswordService() // get password service
	jwtService := jwt_services.NewJwtService()                // get jwt service

	repository := user_repository_implementation.NewUserRepositry(mysql_configs.DB)                // get repository
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService) // get usecase
	delivery := user_delivery_implementation.NewUserDelivery(usecase)                              // get delivery

	routes := r.Group("/user")
	{
		routes.POST("/register", delivery.Register)                                    // register
		routes.POST("/login", delivery.Login)                                          // login
		routes.GET("", middlewares.JWTMiddleware(role.Admin), delivery.GetAllUserData) // fetch all users data
	}
}

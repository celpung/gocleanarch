package user_router

import (
	"net/http"

	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	"github.com/celpung/gocleanarch/configs/role"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/http/user_delivery/implementation"
	user_middleware "github.com/celpung/gocleanarch/delivery/http/user_delivery/middleware"
	user_repository_implementation "github.com/celpung/gocleanarch/domain/user/repository/implementation"
	user_usecase_implementation "github.com/celpung/gocleanarch/domain/user/usecase/implementation"
	jwt_services "github.com/celpung/gocleanarch/services/jwt"
	password_services "github.com/celpung/gocleanarch/services/password"
)

func Router() {
	passwordService := password_services.NewPasswordService()
	jwtService := jwt_services.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql_configs.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	http.HandleFunc("/users/register", delivery.Register)
	http.HandleFunc("/users/login", delivery.Login)
	http.HandleFunc("/users", user_middleware.JWTMiddleware(role.Admin, delivery.GetAllUserData))
	http.HandleFunc("/users/update", user_middleware.JWTMiddleware(role.User, delivery.UpdateUser))
}

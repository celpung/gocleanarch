package user_router

import (
	"net/http"

	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	"github.com/celpung/gocleanarch/configs/role"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/http/user_delivery/implementation"
	middlewares "github.com/celpung/gocleanarch/delivery/http/user_delivery/middleware"
	user_repository_implementation "github.com/celpung/gocleanarch/domain/user/repository/implementation"
	user_usecase_implementation "github.com/celpung/gocleanarch/domain/user/usecase/implementation"
	"github.com/celpung/gocleanarch/services"
	"github.com/celpung/gocleanarch/utils"
)

func Router() {
	passwordService := services.NewPasswordService()
	jwtService := services.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql_configs.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	http.HandleFunc("/users/register", utils.MethodHandler(http.MethodPost, delivery.Register))
	http.HandleFunc("/users/login", utils.MethodHandler(http.MethodPost, delivery.Login))
	http.HandleFunc("/users", utils.MethodHandler(http.MethodGet, middlewares.AuthMiddleware(role.Admin, delivery.GetAllUserData)))
	http.HandleFunc("/users/update", utils.MethodHandler(http.MethodPatch, middlewares.AuthMiddleware(role.User, delivery.UpdateUser)))
	http.HandleFunc("/users/delete", utils.MethodHandler(http.MethodDelete, middlewares.AuthMiddleware(role.Admin, delivery.DeleteUser)))
}

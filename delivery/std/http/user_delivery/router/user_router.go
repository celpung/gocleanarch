package user_router

import (
	"net/http"

	user_repository_implementation "github.com/celpung/gocleanarch/application/user/repository"
	user_usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/std/http/user_delivery/implementation"
	middlewares "github.com/celpung/gocleanarch/delivery/std/http/user_delivery/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
)

func Router() {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	http.HandleFunc("/users/register", middlewares.MethodHandler(http.MethodPost, delivery.Register))
	http.HandleFunc("/users/login", middlewares.MethodHandler(http.MethodPost, delivery.Login))
	http.HandleFunc("/users", middlewares.MethodHandler(http.MethodGet, middlewares.AuthMiddleware(role.Admin, delivery.GetAllUserData)))
	http.HandleFunc("/users/update", middlewares.MethodHandler(http.MethodPatch, middlewares.AuthMiddleware(role.User, delivery.UpdateUser)))
	http.HandleFunc("/users/delete", middlewares.MethodHandler(http.MethodDelete, middlewares.AuthMiddleware(role.Admin, delivery.DeleteUser)))
}

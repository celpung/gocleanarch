package user_router

import (
	"net/http"

	repository_implementation "github.com/celpung/gocleanarch/application/user/repository_implementation"
	usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase_implementation"
	delivery_implementation "github.com/celpung/gocleanarch/delivery/std/http/user/implementation"
	"github.com/celpung/gocleanarch/delivery/std/http/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
)

func Router() {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()

	repository := repository_implementation.NewUserRepository(mysql.DB)
	usecase := usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_implementation.NewUserDelivery(usecase)

	http.HandleFunc("/users/register", middleware.MethodHandler(http.MethodPost, delivery.Register))
	http.HandleFunc("/users/login", middleware.MethodHandler(http.MethodPost, delivery.Login))
	http.HandleFunc("/users", middleware.MethodHandler(http.MethodGet, middleware.AuthMiddleware(role.Admin, delivery.GetAllUserData)))
	http.HandleFunc("/users/update", middleware.MethodHandler(http.MethodPatch, middleware.AuthMiddleware(role.User, delivery.UpdateUser)))
	http.HandleFunc("/users/delete", middleware.MethodHandler(http.MethodDelete, middleware.AuthMiddleware(role.Admin, delivery.DeleteUser)))
}

package user_router

import (
	"net/http"

	repository_impl "github.com/celpung/gocleanarch/application/user/impl/repository"
	usecase_impl "github.com/celpung/gocleanarch/application/user/impl/usecase"
	delivery_impl "github.com/celpung/gocleanarch/delivery/std/http/user/impl"
	"github.com/celpung/gocleanarch/delivery/std/http/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
)

func Router() {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()

	repository := repository_impl.NewUserRepository(mysql.DB)
	usecase := usecase_impl.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_impl.NewUserDelivery(usecase)

	http.HandleFunc("/users/register", middleware.MethodHandler(http.MethodPost, delivery.Register))
	http.HandleFunc("/users/login", middleware.MethodHandler(http.MethodPost, delivery.Login))
	http.HandleFunc("/users", middleware.MethodHandler(http.MethodGet, middleware.AuthMiddleware(role.Admin, delivery.GetAllUserData)))
	http.HandleFunc("/search", middleware.MethodHandler(http.MethodGet, middleware.AuthMiddleware(role.Admin, delivery.SearchUser)))
	http.HandleFunc("/users/update", middleware.MethodHandler(http.MethodPatch, middleware.AuthMiddleware(role.User, delivery.UpdateUser)))
	http.HandleFunc("/users/delete", middleware.MethodHandler(http.MethodDelete, middleware.AuthMiddleware(role.Admin, delivery.DeleteUser)))
}

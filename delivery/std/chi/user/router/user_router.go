package user_router

import (
	"github.com/go-chi/chi/v5"

	repository_impl "github.com/celpung/gocleanarch/application/user/impl/repository"
	usecase_impl "github.com/celpung/gocleanarch/application/user/impl/usecase"
	delivery_impl "github.com/celpung/gocleanarch/delivery/std/chi/user/impl"
	"github.com/celpung/gocleanarch/delivery/std/chi/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
)

// Router mendaftarkan semua route user ke router utama
func Router(r chi.Router) {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()

	repository := repository_impl.NewUserRepository(mysql.DB)
	usecase := usecase_impl.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_impl.NewUserDelivery(usecase)

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", delivery.Register)
		r.Post("/login", delivery.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(middleware.Admin, middleware.Super))
			r.Get("/", delivery.GetAllUserData)
			r.Delete("/{id}", delivery.DeleteUser)
			r.Get("/search", delivery.SearchUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(middleware.User, middleware.Admin, middleware.Super))
			r.Patch("/update", delivery.UpdateUser)
		})
	})
}

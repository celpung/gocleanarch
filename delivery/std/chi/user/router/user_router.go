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
		// Public routes
		r.Post("/register", delivery.Register)
		r.Post("/login", delivery.Login)

		// Protected routes (Admin and super)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(middleware.Admin, middleware.Super))
			r.Get("/", delivery.GetAllUserData)
			r.Delete("/delete", delivery.DeleteUser)
			r.Delete("/search", delivery.SearchUser)
		})

		// Protected routes (User)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(middleware.User))
			r.Patch("/update", delivery.UpdateUser)
		})
	})
}

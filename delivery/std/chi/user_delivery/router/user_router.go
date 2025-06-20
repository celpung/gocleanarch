package user_router

import (
	"github.com/go-chi/chi/v5"

	user_repository_implementation "github.com/celpung/gocleanarch/application/user/repository"
	user_usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/std/chi/user_delivery/implementation"
	middlewares "github.com/celpung/gocleanarch/delivery/std/chi/user_delivery/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	mysql_configs "github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
)

// Router mendaftarkan semua route user ke router utama
func Router(r chi.Router) {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()

	repository := user_repository_implementation.NewUserRepository(mysql_configs.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	r.Route("/users", func(r chi.Router) {
		// Public routes
		r.Post("/register", delivery.Register)
		r.Post("/login", delivery.Login)

		// Protected routes (Admin)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(role.Admin))
			r.Get("/", delivery.GetAllUserData)
			r.Delete("/delete", delivery.DeleteUser)
		})

		// Protected routes (User)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(role.User))
			r.Patch("/update", delivery.UpdateUser)
		})
	})
}

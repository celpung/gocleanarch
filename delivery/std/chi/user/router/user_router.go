package user_router

import (
	"github.com/go-chi/chi/v5"

	repository_implementation "github.com/celpung/gocleanarch/application/user/repository_implementation"
	usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase_implementation"
	delivery_implementation "github.com/celpung/gocleanarch/delivery/std/chi/user/implementation"
	"github.com/celpung/gocleanarch/delivery/std/chi/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
)

// Router mendaftarkan semua route user ke router utama
func Router(r chi.Router) {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()

	repository := repository_implementation.NewUserRepository(mysql.DB)
	usecase := usecase_implementation.NewUserUsecase(repository, passwordService, jwtService)
	delivery := delivery_implementation.NewUserDelivery(usecase)

	r.Route("/users", func(r chi.Router) {
		// Public routes
		r.Post("/register", delivery.Register)
		r.Post("/login", delivery.Login)

		// Protected routes (Admin)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(role.Admin))
			r.Get("/", delivery.GetAllUserData)
			r.Delete("/delete", delivery.DeleteUser)
		})

		// Protected routes (User)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(role.User))
			r.Patch("/update", delivery.UpdateUser)
		})
	})
}

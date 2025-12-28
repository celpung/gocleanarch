package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/delivery/handler"
	"github.com/celpung/gocleanarch/internal/delivery/middleware"
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
)

func UserRouter(r chi.Router, verifier dependencies.TokenVerifier, userHandler *handler.UserHandler) {
	r.Route("/user", func(r chi.Router) {
		// PUBLIC
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		// PROTECTED: User/Admin/Super
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(verifier, middleware.User, middleware.Admin, middleware.Super))
			r.Post("/change-password", userHandler.ChangePassword)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// PROTECTED: Admin/Super
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(verifier, middleware.Admin, middleware.Super))
			r.Get("/", userHandler.ListUsers)
		})
	})
}

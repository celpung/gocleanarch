package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/delivery/handler"
	"github.com/celpung/gocleanarch/internal/delivery/middleware"
)

func UserRouter(r chi.Router, userHandler *handler.UserHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Get("/login", userHandler.Login)

		r.Use(middleware.AuthMiddleware(middleware.User, middleware.Admin, middleware.Super))
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)

		r.Use(middleware.AuthMiddleware(middleware.Admin, middleware.Super))
		r.Get("/", userHandler.ListUsers)
	})
}

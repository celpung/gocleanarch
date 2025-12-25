package user_router

import (
	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal/modules/user/delivery/handler"
)

func Register(r chi.Router, userHandler *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/login", userHandler.Login)
		r.Get("/", userHandler.ListUsers)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}

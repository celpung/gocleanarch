package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/celpung/gocleanarch/internal_old/delivery/handler"
)

func UserRouter(r chi.Router, userHandler *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.ListUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}

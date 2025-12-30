package router

import (
	"github.com/celpung/gocleanarch/internal/delivery/handler"
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, handler *handler.UserHandler) {
	r.Post("/register", handler.Register)
}

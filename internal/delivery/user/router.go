package user

import (
	"github.com/go-chi/chi/v5"

	usecase "github.com/celpung/gocleanarch/internal/usecase/user"
)

func RegisterRoutes(r chi.Router, verifier usecase.TokenVerifier, handler *Handler) {
	r.Route("/user", func(r chi.Router) {
		// PUBLIC
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		// PROTECTED: User/Admin/Super
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(verifier, User, Admin, Super))
			r.Post("/change-password", handler.ChangePassword)
			r.Put("/{id}", handler.UpdateUser)
			r.Delete("/{id}", handler.DeleteUser)
		})

		// PROTECTED: Admin/Super
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(verifier, Admin, Super))
			r.Get("/", handler.ListUsers)
		})
	})
}

package user_router

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/celpung/gocleanarch/app/modules/user/handler"
	"github.com/celpung/gocleanarch/app/modules/user/repository"
	"github.com/celpung/gocleanarch/app/modules/user/usecase"
)

/*
Router adalah composition root untuk modul user.
Semua dependency DIRAKIT di sini, tapi DITERIMA dari luar.
*/

func Register(r chi.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
	})
}

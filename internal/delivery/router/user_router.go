package router

// jwtSvc := services.NewJwtService()
// secret := environment.Env.JWT_SECRET // dibaca di layer infra/config

// r.Route("/users", func(r chi.Router) {
// 	r.Use(middleware.AuthMiddleware(jwtSvc, secret, middleware.Admin, middleware.Super))
// 	r.Get("/", userHandler.Lists)

// 	// endpoint yang hanya butuh login (tanpa role guard)
// 	r.With(middleware.AuthMiddleware(jwtSvc, secret)).Get("/me", userHandler.Me)
// })

// func RegisterRoutes(r chi.Router, verifier usecase.TokenVerifier, handler *Handler) {
// 	r.Route("/user", func(r chi.Router) {
// 		// PUBLIC
// 		r.Post("/register", handler.Register)
// 		r.Post("/login", handler.Login)

// 		// PROTECTED: User/Admin/Super
// 		r.Group(func(r chi.Router) {
// 			r.Use(AuthMiddleware(verifier, User, Admin, Super))
// 			r.Post("/change-password", handler.ChangePassword)
// 			r.Put("/{id}", handler.UpdateUser)
// 			r.Delete("/{id}", handler.DeleteUser)
// 		})

// 		// PROTECTED: Admin/Super
// 		r.Group(func(r chi.Router) {
// 			r.Use(AuthMiddleware(verifier, Admin, Super))
// 			r.Get("/", handler.ListUsers)
// 		})
// 	})
// }

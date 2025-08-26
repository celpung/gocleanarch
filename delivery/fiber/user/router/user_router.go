package user_router

import (
	repository_implementation "github.com/celpung/gocleanarch/application/user/repository_implementation"
	usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase_implementation"
	delivery_implementation "github.com/celpung/gocleanarch/delivery/fiber/user/implementation"
	middleware "github.com/celpung/gocleanarch/delivery/fiber/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRouter(router fiber.Router) {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()
	repo := repository_implementation.NewUserRepository(mysql.DB)
	usecase := usecase_implementation.NewUserUsecase(repo, passwordService, jwtService)
	delivery := delivery_implementation.NewUserDelivery(usecase)

	user := router.Group("/users")
	user.Post("/register", delivery.Register)
	user.Post("/login", delivery.Login)
	user.Get("/", middleware.AuthMiddleware(role.Admin), delivery.GetAllUserData)
	user.Patch("/", middleware.AuthMiddleware(role.Admin), delivery.UpdateUser)
	user.Delete("/:user_id", middleware.AuthMiddleware(role.Admin), delivery.DeleteUser)
}

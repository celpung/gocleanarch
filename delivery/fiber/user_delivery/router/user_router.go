package user_router

import (
	user_repository_implementation "github.com/celpung/gocleanarch/application/user/repository"
	user_usecase_implementation "github.com/celpung/gocleanarch/application/user/usecase"
	user_delivery_implementation "github.com/celpung/gocleanarch/delivery/fiber/user_delivery/implementation"
	middlewares "github.com/celpung/gocleanarch/delivery/fiber/user_delivery/middlewara"
	"github.com/celpung/gocleanarch/infrastructure/auths"
	mysql_configs "github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRouter(router fiber.Router) {
	passwordService := auths.NewPasswordService()
	jwtService := auths.NewJwtService()
	repo := user_repository_implementation.NewUserRepository(mysql_configs.DB)
	usecase := user_usecase_implementation.NewUserUsecase(repo, passwordService, jwtService)
	delivery := user_delivery_implementation.NewUserDelivery(usecase)

	user := router.Group("/users")
	user.Post("/register", delivery.Register)
	user.Post("/login", delivery.Login)
	user.Get("/", middlewares.AuthMiddleware(role.Admin), delivery.GetAllUserData)
	user.Patch("/", middlewares.AuthMiddleware(role.Admin), delivery.UpdateUser)
	user.Delete("/:user_id", middlewares.AuthMiddleware(role.Admin), delivery.DeleteUser)
}

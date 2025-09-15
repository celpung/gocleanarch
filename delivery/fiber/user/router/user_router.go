package user_router

import (
	repository_impl "github.com/celpung/gocleanarch/application/user/impl/repository"
	usecase_impl "github.com/celpung/gocleanarch/application/user/impl/usecase"
	delivery_impl "github.com/celpung/gocleanarch/delivery/fiber/user/impl"
	middleware "github.com/celpung/gocleanarch/delivery/fiber/user/middleware"
	"github.com/celpung/gocleanarch/infrastructure/auth"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRouter(router fiber.Router) {
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJwtService()
	repo := repository_impl.NewUserRepository(mysql.DB)
	usecase := usecase_impl.NewUserUsecase(repo, passwordService, jwtService)
	delivery := delivery_impl.NewUserDelivery(usecase)

	user := router.Group("/users")
	user.Post("/register", delivery.Register)
	user.Post("/login", delivery.Login)
	user.Get("/", middleware.AuthMiddleware(role.Admin), delivery.GetAllUserData)
	user.Get("/search", middleware.AuthMiddleware(role.Admin), delivery.SearchUser)
	user.Patch("/", middleware.AuthMiddleware(role.Admin), delivery.UpdateUser)
	user.Delete("/:user_id", middleware.AuthMiddleware(role.Admin), delivery.DeleteUser)
}

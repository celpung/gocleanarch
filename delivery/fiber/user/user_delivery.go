package delivery

import "github.com/gofiber/fiber/v2"

type UserDelivery interface {
	Register(c *fiber.Ctx) error
	GetAllUserData(c *fiber.Ctx) error
	SearchUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

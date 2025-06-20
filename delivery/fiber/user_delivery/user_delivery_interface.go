package user_delivery

import "github.com/gofiber/fiber/v2"

type UserDeliveryInterface interface {
	Register(c *fiber.Ctx) error
	GetAllUserData(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

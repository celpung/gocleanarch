package main

import (
	"log"

	user_router "github.com/celpung/gocleanarch/delivery/fiber/user_delivery/router"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/environment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to the database and auto migrate
	mysql.CreateDatabaseIfNotExists()
	mysql.ConnectDatabase()
	mysql.AutoMigrage()

	// setup mode
	mode := environment.Env.MODE

	r := fiber.New(fiber.Config{
		AppName:               "Skoolar Auth",
		DisableStartupMessage: mode == "release",
	})

	allowedOrigins := environment.Env.ALLOWED_ORIGINS

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	if mode == "debug" {
		r.Use(func(c *fiber.Ctx) error {
			log.Printf("[DEBUG] %s %s", c.Method(), c.Path())
			return c.Next()
		})
	}

	api := r.Group("/api")
	user_router.RegisterUserRouter(api)

	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("../../public/index.html")
	})

	r.Static("/images", "../../public/images")

	log.Printf("Running in %s mode", mode)
	log.Fatal(r.Listen(":" + environment.Env.PORT))
}

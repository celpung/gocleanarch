package main

import (
	"log"

	user_router "github.com/celpung/gocleanarch/delivery/fiber/user_delivery/router"
	mysql_configs "github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/environment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to the database and auto migrate
	mysql_configs.CreateDatabaseIfNotExists()
	mysql_configs.ConnectDatabase()
	mysql_configs.AutoMigrage()

	// setup mode
	mode := environment.Env.MODE

	app := fiber.New(fiber.Config{
		AppName:               "Skoolar Auth",
		DisableStartupMessage: mode == "release",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	if mode == "debug" {
		app.Use(func(c *fiber.Ctx) error {
			log.Printf("[DEBUG] %s %s", c.Method(), c.Path())
			return c.Next()
		})
	}

	api := app.Group("/api")
	user_router.RegisterUserRouter(api)

	log.Printf("Running in %s mode", mode)
	log.Fatal(app.Listen(":" + environment.Env.PORT))
}

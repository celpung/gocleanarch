package main

import (
	"fmt"
	"log"
	"os"

	"github.com/celpung/gocleanarch/configs"
	"github.com/celpung/gocleanarch/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env", err)
	}

	// Connect to the database
	configs.ConnectDatabase()
	// do the automigrate
	configs.AutoMigrage()

	//setup gin
	r := gin.Default()

	// setup mode
	mode := os.Getenv("MODE")

	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Please set the mode debug/release on environment!")
		fmt.Println("Example : [MODE: debug] or [MODE: release]")
		fmt.Println("-------------------------------------------------")
		panic("Critical error, cannot find mode on environment!")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// import all routers
	user.Router(r)

	// Start the server
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

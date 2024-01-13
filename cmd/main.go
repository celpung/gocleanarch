package main

import (
	"fmt"
	"log"
	"os"

	"github.com/celpung/gocleanarch/configs"
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
	configs.AutoMigrage()

	//setup gin
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// import all routers
	// routers.UserRouter(r)

	// Start the server
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

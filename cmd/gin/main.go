package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	"github.com/celpung/gocleanarch/configs/environment"
	user_router "github.com/celpung/gocleanarch/domain/user/delivery/gin/router"
	"github.com/celpung/gocleanarch/entity"
	crud_router "github.com/celpung/gocleanarch/utils/crud/delivery/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	// Connect to the database and auto migrate
	mysql_configs.CreateDatabaseIfNotExists()
	mysql_configs.ConnectDatabase()
	mysql_configs.AutoMigrage()

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
		// panic("Critical error, cannot find mode on environment!")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:  environment.Env.AllowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "x-api-key", "X-API-Key"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// setup router
	api := r.Group("/api")
	user_router.Router(api)

	crud_router.SetupRouter[entity.Slider](
		api,
		mysql_configs.DB,
		reflect.TypeOf(entity.Slider{}),
		"/sliders",
		map[string][]gin.HandlerFunc{
			"POST":   {},
			"READ":   {},
			"PUT":    {},
			"DELETE": {},
		})

	// Serve static files
	r.GET("/", func(c *gin.Context) {
		c.File("../../public/index.html")
	})

	r.Static("/images", "../../public/images")

	// Start the server
	r.Run(fmt.Sprintf(":%s", environment.Env.Port))
}

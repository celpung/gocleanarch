package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	crud_router "github.com/celpung/go-generic-crud/crud_router"
	mysql_configs "github.com/celpung/gocleanarch/configs/database/mysql"
	"github.com/celpung/gocleanarch/configs/environment"
	user_router "github.com/celpung/gocleanarch/domain/user/delivery/gin/router"
	"github.com/celpung/gocleanarch/entity"
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

	// setup mode
	mode := environment.Env.MODE

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	//setup gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
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
	r.Run(fmt.Sprintf(":%s", environment.Env.PORT))
}

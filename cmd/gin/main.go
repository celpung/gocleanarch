package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	crud_router "github.com/celpung/go-generic-crud/crud_router"
	user_router "github.com/celpung/gocleanarch/delivery/gin/user/router"
	slider_entity "github.com/celpung/gocleanarch/domain/slider/entity"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database and auto migrate
	mysql.CreateDatabaseIfNotExists()
	mysql.ConnectDatabase()
	mysql.AutoMigrate()

	// setup mode
	mode := environment.Env.MODE

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	//setup gin
	r := gin.Default()

	allowedOrigins := strings.Split(environment.Env.ALLOWED_ORIGINS, ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// setup router
	api := r.Group("/api")
	user_router.Router(api)

	// implement generic CRUD router
	crud_router.SetupRouter[slider_entity.Slider](
		api,
		mysql.DB,
		reflect.TypeOf(slider_entity.Slider{}),
		"/sliders",
		map[string][]gin.HandlerFunc{
			"POST":   {},
			"GET":    {},
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

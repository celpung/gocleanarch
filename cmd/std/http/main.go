package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	user_router "github.com/celpung/gocleanarch/delivery/std/http/user/router"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/environment"
)

func main() {
	// Connect to the database and auto migrate
	mysql.CreateDatabaseIfNotExists()
	mysql.ConnectDatabase()
	mysql.AutoMigrage()

	// Setup mode
	mode := environment.Env.MODE

	if mode != "debug" && mode != "release" {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Please set the mode debug/release in the environment!")
		fmt.Println("Example: [MODE=debug] or [MODE=release]")
		fmt.Println("-------------------------------------------------")
		panic("Critical error: cannot find mode in environment!")
	}

	// Setup CORS headers
	allowedOrigins := environment.Env.ALLOWED_ORIGINS
	if allowedOrigins == "" {
		log.Fatal("ALLOWED_ORIGINS environment variable is not set")
	}
	origins := strings.Split(allowedOrigins, ",")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if r.Method == http.MethodOptions {
			return
		}

		http.ServeFile(w, r, "../../public/index.html")
	})

	user_router.Router()

	// Serve static files
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("../../public/images"))))

	// Log the startup message in debug mode
	port := environment.Env.PORT
	if mode == "debug" {
		log.Printf("Application is running in debug mode on port %s", port)
	}

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

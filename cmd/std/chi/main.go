package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	user_router "github.com/celpung/gocleanarch/delivery/std/chi/user/router"
	"github.com/celpung/gocleanarch/infrastructure/db/mysql"
	"github.com/celpung/gocleanarch/infrastructure/environment"
)

func main() {
	// Connect to the database and auto migrate
	if err := mysql.CreateDatabaseIfNotExists(); err != nil {
		log.Fatalf("failed to prepare database: %v", err)
	}
	if err := mysql.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := mysql.AutoMigrate(); err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}

	// Setup mode
	mode := environment.Env.MODE
	if mode != "debug" && mode != "release" {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Please set the mode debug/release in the environment!")
		fmt.Println("Example: [MODE=debug] or [MODE=release]")
		fmt.Println("-------------------------------------------------")
		panic("Critical error: cannot find mode in environment!")
	}

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Get allowed origins from environment
	allowedOriginsRaw := environment.Env.ALLOWED_ORIGINS
	if allowedOriginsRaw == "" {
		log.Fatal("ALLOWED_ORIGINS environment variable is not set")
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	// CORS Middleware using environment
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(origin) == strings.TrimSpace(allowed) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Static index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../public/index.html")
	})

	// Static images
	fileServer := http.StripPrefix("/images", http.FileServer(http.Dir("../../public/images")))
	r.Handle("/images/*", fileServer)

	// Register user routes
	user_router.Router(r)

	// Start server
	port := environment.Env.PORT
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start chi server: %v", err)
	}
}

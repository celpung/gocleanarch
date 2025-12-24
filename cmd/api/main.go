package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/app/infra/db/mysql"
	"github.com/celpung/gocleanarch/app/infra/environment"
	user_router "github.com/celpung/gocleanarch/app/modules/user/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// VALIDATE MODE
	mode := environment.Env.MODE
	if mode != "debug" && mode != "release" {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Please set MODE=debug or MODE=release")
		fmt.Println("-------------------------------------------------")
		panic("invalid MODE environment")
	}

	// DATABASE BOOTSTRAP
	dbCfg := mysql.Config{
		Username: environment.Env.DB_USERNAME,
		Password: environment.Env.DB_PASSWORD,
		Host:     environment.Env.DB_HOST,
		Port:     environment.Env.DB_PORT,
		Database: environment.Env.DB_NAME,
	}

	database, err := mysql.New(dbCfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	db := database.DB

	// ROUTER
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	allowedOriginsRaw := environment.Env.ALLOWED_ORIGINS
	if allowedOriginsRaw == "" {
		log.Fatal("ALLOWED_ORIGINS environment variable is not set")
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	r.Use(corsMiddleware(allowedOrigins))

	// STATIC FILES
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	fileServer := http.StripPrefix(
		"/images",
		http.FileServer(http.Dir("./public/images")),
	)
	r.Handle("/images/*", fileServer)

	// MODULE ROUTES
	user_router.Register(r, db)

	// START SERVER
	port := environment.Env.PORT
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// CORS MIDDLEWARE
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
	}
}

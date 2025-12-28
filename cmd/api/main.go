package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/internal/delivery/handler"
	"github.com/celpung/gocleanarch/internal/delivery/router"
	"github.com/celpung/gocleanarch/internal/infra/auth"
	"github.com/celpung/gocleanarch/internal/infra/db/mysql"
	"github.com/celpung/gocleanarch/internal/infra/environment"
	"github.com/celpung/gocleanarch/internal/infra/identity"
	"github.com/celpung/gocleanarch/internal/infra/persistence"
	"github.com/celpung/gocleanarch/internal/usecase"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	env := environment.Load()

	// VALIDATE MODE
	mode := env.MODE
	if mode != "debug" && mode != "release" {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Please set MODE=debug or MODE=release")
		fmt.Println("-------------------------------------------------")
		panic("invalid MODE environment")
	}

	// DATABASE BOOTSTRAP
	dbCfg := mysql.Config{
		Username: env.DB_USERNAME,
		Password: env.DB_PASSWORD,
		Host:     env.DB_HOST,
		Port:     env.DB_PORT,
		Database: env.DB_NAME,
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
	allowedOriginsRaw := env.ALLOWED_ORIGINS
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

	userRepo := persistence.NewUserRepository(db)
	jwtGenerator := auth.NewJwtGenerator(env.JWT_TOKEN)
	jwtVerifier := auth.NewJwtVerifier(env.JWT_TOKEN)
	userUsecase := usecase.NewUserUsecase(
		userRepo,
		identity.UUIDGenerator{},
		auth.BcryptHasher{},
		jwtGenerator,
		typograph.Typograph{},
	)
	userHandler := handler.NewUserHandler(userUsecase)

	router.UserRouter(r, jwtVerifier, userHandler)

	// START SERVER
	port := env.PORT
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

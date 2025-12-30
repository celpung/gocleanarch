package container

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	deliveryuser "github.com/celpung/gocleanarch/internal_old/delivery/user"
	"github.com/celpung/gocleanarch/internal_old/infra/db"
	"github.com/celpung/gocleanarch/internal_old/infra/environment"
	"github.com/celpung/gocleanarch/internal_old/infra/persistence"
	"github.com/celpung/gocleanarch/internal_old/infra/services"
	usecase "github.com/celpung/gocleanarch/internal_old/usecase/user"
	"github.com/celpung/gocleanarch/pkg/typograph"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

// Container wires application dependencies and exposes the HTTP router.
type Container struct {
	Env         environment.Environment
	DB          *gorm.DB
	Router      http.Handler
	UserUsecase usecase.Usecase
}

func Build() (*Container, error) {
	env := environment.Load()
	if err := validateMode(env.MODE); err != nil {
		return nil, err
	}

	dbProvider := db.DefaultProvider()
	dbConn, err := dbProvider.Connect(env)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	userRepo := persistence.NewUserRepository(dbConn)
	jwtService := services.NewJWTService(env.JWT_TOKEN)
	userUsecase := usecase.NewUsecase(
		userRepo,
		services.UUIDService{},
		services.PasswordService{},
		jwtService,
		typograph.Typograph{},
	)

	r, err := buildRouter(env, jwtService, userUsecase)
	if err != nil {
		return nil, err
	}

	return &Container{
		Env:         env,
		DB:          dbConn,
		Router:      r,
		UserUsecase: userUsecase,
	}, nil
}

func buildRouter(env environment.Environment, verifier usecase.TokenVerifier, userUsecase usecase.Usecase) (http.Handler, error) {
	allowedOriginsRaw := env.ALLOWED_ORIGINS
	if allowedOriginsRaw == "" {
		return nil, fmt.Errorf("ALLOWED_ORIGINS environment variable is not set")
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
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

	userHandler := deliveryuser.NewHandler(userUsecase)
	deliveryuser.RegisterRoutes(r, verifier, userHandler)
	return r, nil
}

func validateMode(mode string) error {
	if mode != "debug" && mode != "release" {
		return fmt.Errorf("please set MODE=debug or MODE=release")
	}
	return nil
}

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

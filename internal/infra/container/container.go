package container

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celpung/gocleanarch/internal/delivery/handler"
	"github.com/celpung/gocleanarch/internal/delivery/router"
	"github.com/celpung/gocleanarch/internal/infra/db"
	"github.com/celpung/gocleanarch/internal/infra/db/migration"
	"github.com/celpung/gocleanarch/internal/infra/environment"
	"github.com/celpung/gocleanarch/internal/infra/repository"
	"github.com/celpung/gocleanarch/internal/usecase"
	usecaseport "github.com/celpung/gocleanarch/internal/usecase/port"
	"github.com/celpung/gocleanarch/pkg/services"
	"github.com/celpung/gocleanarch/pkg/utilities/typograph"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Container struct {
	Env         environment.Environment
	DB          *gorm.DB
	Router      http.Handler
	UserUsecase usecaseport.UserUsecase
}

func Build() (*Container, error) {
	// load env
	env := environment.Load()

	// validate mode
	if err := validateMode(env.MODE); err != nil {
		return nil, err
	}

	// setup database
	dbProvider := db.DefaultProvider()
	dbConn, err := dbProvider.Connect(env)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := migration.Run(dbConn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	deps := dependencies{
		uuid:      services.UUIDService{},
		password:  services.PasswordService{},
		typograph: typograph.Typograph{},
		jwt:       services.NewJwtService(),
		secret:    env.JWT_TOKEN,
	}

	// register modules
	userRepo := repository.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo, deps)
	userHandler := handler.NewUserHandler(userUsecase)

	r := chi.NewRouter()
	router.Router(r, userHandler)

	return &Container{
		Env:         env,
		DB:          dbConn,
		Router:      r,
		UserUsecase: userUsecase,
	}, nil

}

func validateMode(mode string) error {
	if mode != "debug" && mode != "release" {
		return fmt.Errorf("please set MODE=debug or MODE=release")
	}
	return nil
}

type dependencies struct {
	uuid      services.UUIDService
	password  services.PasswordService
	typograph typograph.Typograph
	jwt       *services.JwtService
	secret    string
}

func (d dependencies) NewID() (string, error) {
	return d.uuid.NewID()
}

func (d dependencies) HashPassword(plain string) (string, error) {
	return d.password.Hash(plain)
}

func (d dependencies) ComparePassword(hashed string, plain string) error {
	return d.password.Compare(hashed, plain)
}

func (d dependencies) GenerateToken(id string, email string, role string) (string, error) {
	return d.jwt.GenerateToken(id, email, role, d.secret, 24*time.Hour)
}

func (d dependencies) ToTitleCase(s string) string {
	return d.typograph.ToTitleCase(s)
}

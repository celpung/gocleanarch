package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	BASE_URL        string
	PORT            string
	APP_NAME        string
	MODE            string
	JWT_SECRET      string
	DB_USERNAME     string
	DB_PASSWORD     string
	DB_NAME         string
	DB_PORT         string
	DB_HOST         string
	DB_DIALECT      string
	ALLOWED_ORIGINS string
}

var Env Environment

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using hardcoded values")
	}

	// Initialize environment variables with fallback to hardcoded defaults
	Env = Environment{
		BASE_URL:        getEnv("BASE_URL", "http://localhost"),
		PORT:            getEnv("PORT", "8080"),
		APP_NAME:        getEnv("APP_NAME", "GoCleanArch"),
		MODE:            getEnv("MODE", "debug"),
		JWT_SECRET:      getEnv("MODE", "534LK786HJK7DHFG89"),
		DB_USERNAME:     getEnv("DB_USERNAME", "root"),
		DB_PASSWORD:     getEnv("DB_PASSWORD", ""),
		DB_NAME:         getEnv("DB_NAME", "gocleanarch"),
		DB_PORT:         getEnv("DB_PORT", "3306"),
		DB_HOST:         getEnv("DB_HOST", "127.0.0.1"),
		DB_DIALECT:      getEnv("DB_DIALECT", "mysql"),
		ALLOWED_ORIGINS: getEnv("ALLOWED_ORIGINS", "http://localhost"),
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

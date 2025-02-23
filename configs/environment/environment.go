package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	BASE_URL               string
	PORT                   string
	EMAIL_CONFIRMATION_URL string
	APP_NAME               string
	MODE                   string
}

var Env Environment

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using hardcoded values")
	}

	// Initialize environment variables with fallback to hardcoded defaults
	Env = Environment{
		BASE_URL:               getEnv("BASE_URL", "http://localhost"),
		PORT:                   getEnv("PORT", "8080"),
		EMAIL_CONFIRMATION_URL: getEnv("EMAIL_CONFIRMATION_URL", "http://localhost:8080/api/users/verify-email"),
		APP_NAME:               getEnv("APP_NAME", "GoCleanArch"),
		MODE:                   getEnv("MODE", "debug"),
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

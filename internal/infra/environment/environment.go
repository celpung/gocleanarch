package environment

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Environment struct {
	BASE_URL        string
	PORT            string
	APP_NAME        string
	MODE            string
	JWT_TOKEN       string
	DB_USERNAME     string
	DB_PASSWORD     string
	DB_NAME         string
	DB_PORT         string
	DB_HOST         string
	DB_DIALECT      string
	ALLOWED_ORIGINS string
	AUTO_MIGRATE    bool
}

// Load reads environment variables (or .env) and returns a configured struct.
func Load() Environment {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using hardcoded values")
	}

	// Initialize environment variables with fallback to hardcoded defaults
	return Environment{
		BASE_URL:        getEnv("BASE_URL", "http://localhost"),
		PORT:            getEnv("PORT", "8080"),
		APP_NAME:        getEnv("APP_NAME", "GoCleanArch"),
		MODE:            getEnv("MODE", "debug"),
		JWT_TOKEN:       getEnv("JWT_TOKEN", "534LK786HJK7DHFG89"),
		DB_USERNAME:     getEnv("DB_USERNAME", "root"),
		DB_PASSWORD:     getEnv("DB_PASSWORD", ""),
		DB_NAME:         getEnv("DB_NAME", "gocleanarch"),
		DB_PORT:         getEnv("DB_PORT", "3306"),
		DB_HOST:         getEnv("DB_HOST", "127.0.0.1"),
		DB_DIALECT:      getEnv("DB_DIALECT", "mysql"),
		ALLOWED_ORIGINS: getEnv("ALLOWED_ORIGINS", "http://localhost,http://localhost:5173,http://localhost:3000"),
		AUTO_MIGRATE:    getEnvAsBool("AUTO_MIGRATE", true),
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return parsed
}

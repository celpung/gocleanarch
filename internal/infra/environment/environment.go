package environment

import (
	"fmt"
	"log"
	"os"
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
}

// Load reads environment variables (or .env) and returns a configured struct.
func Load() (Environment, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using hardcoded values")
	}

	// Initialize environment variables with fallback to hardcoded defaults
	env := Environment{
		BASE_URL:        strings.TrimSpace(getEnv("BASE_URL", "http://localhost")),
		PORT:            strings.TrimSpace(getEnv("PORT", "8080")),
		APP_NAME:        strings.TrimSpace(getEnv("APP_NAME", "GoCleanArch")),
		MODE:            strings.TrimSpace(getEnv("MODE", "debug")),
		JWT_TOKEN:       strings.TrimSpace(getEnv("JWT_TOKEN", "534LK786HJK7DHFG89")),
		DB_USERNAME:     strings.TrimSpace(getEnv("DB_USERNAME", "root")),
		DB_PASSWORD:     strings.TrimSpace(getEnv("DB_PASSWORD", "")),
		DB_NAME:         strings.TrimSpace(getEnv("DB_NAME", "gocleanarch")),
		DB_PORT:         strings.TrimSpace(getEnv("DB_PORT", "3306")),
		DB_HOST:         strings.TrimSpace(getEnv("DB_HOST", "127.0.0.1")),
		DB_DIALECT:      strings.TrimSpace(getEnv("DB_DIALECT", "mysql")),
		ALLOWED_ORIGINS: strings.TrimSpace(getEnv("ALLOWED_ORIGINS", "http://localhost,http://localhost:5173,http://localhost:3000")),
	}

	if err := validateEnv(env); err != nil {
		return Environment{}, err
	}

	return env, nil
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func validateEnv(env Environment) error {
	if strings.TrimSpace(env.JWT_TOKEN) == "" {
		return fmt.Errorf("JWT_TOKEN is required")
	}

	if !strings.EqualFold(env.DB_DIALECT, "mysql") {
		return fmt.Errorf("unsupported db dialect: %s (only mysql is allowed)", env.DB_DIALECT)
	}

	return nil
}

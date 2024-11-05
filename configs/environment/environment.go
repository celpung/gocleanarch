package environment

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
)

type Environment struct {
	BaseURL              string
	Port                 string
	EmailConfirmationURL string
	AppName              string
	AllowedOrigins       []string
}

var Env Environment

func init() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from hardcoded values ")
	}

	// Initialize the environment variables (either from .env or hardcoded)
	Env = Environment{
		BaseURL:              "http://localhost:8080",
		Port:                 "8080",
		EmailConfirmationURL: "http://localhost:8080",
		AppName:              "Gocleanarch",
		AllowedOrigins:       strings.Split("http://localhost,http://localhost:5173", ","),
	}
}

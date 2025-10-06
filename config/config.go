// 📁 config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// JWTSecret holds the secret key for JWT.
var JWTSecret string

// LoadConfig loads environment variables from .env file.
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("Error: JWT_SECRET environment variable not set")
	}
}
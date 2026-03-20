package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	// Add your configuration fields here
	// Example:
	// ServerPort string
	// LogLevel   string
}

// Load reads environment variables from .env file and returns Config
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{}

	// Load your environment variables here
	// Example:
	// cfg.ServerPort = os.Getenv("SERVER_PORT")
	// if cfg.ServerPort == "" {
	//     cfg.ServerPort = "8080"
	// }

	return cfg
}

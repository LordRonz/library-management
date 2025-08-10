package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds the server configuration.
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load loads the configuration from environment variables.
func Load() *Config {
	// Load .env file
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "library"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
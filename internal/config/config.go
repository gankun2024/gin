package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Stripe   StripeConfig
}

// ServerConfig holds server related configuration
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds database related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// AuthConfig holds authentication related configuration
type AuthConfig struct {
	JWTSecret    string
	TokenExpiry  int
	RefreshToken bool
}

// StripeConfig holds Stripe related configuration
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:    getEnv("JWT_SECRET", ""),
			TokenExpiry:  getEnvAsInt("JWT_EXPIRY", 24),
			RefreshToken: getEnvAsBool("JWT_REFRESH_ENABLED", true),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
	}

	// Validate required configurations
	// if cfg.Auth.JWTSecret == "" {
	// 	return nil, fmt.Errorf("JWT_SECRET is required")
	// }

	// if cfg.Stripe.SecretKey == "" {
	// 	return nil, fmt.Errorf("STRIPE_SECRET_KEY is required")
	// }

	return cfg, nil
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := fmt.Sscanf(value, "%d"); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" {
			return true
		}
		if value == "false" || value == "0" {
			return false
		}
	}
	return defaultValue
}

package wizards

import (
	"os"
)

// Config holds all the configuration for the application.
type Config struct {
	CoreServiceURL string
	ServerPort     string
	// Add other configurations like database DSN, JWT secret, etc. here
}

// NewConfig creates a new Config struct and loads values from environment variables.
func NewConfig() *Config {
	return &Config{
		CoreServiceURL: getEnv("CORE_SERVICE_URL", "http://localhost:8080"), // Default for local dev
		ServerPort:     getEnv("SERVER_PORT", ":8081"),                      // Default for payment-service
	}
}

// Helper function to read an environment variable or return a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

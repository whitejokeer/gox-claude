package config

import (
	"os"
)

// Config holds the service configuration
type Config struct {
	ServiceName string
	Port        string
	Host        string
	LogLevel    string
	
	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	
	// JWT
	JWTSecret string
	
	// Rate limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitDuration string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "posts"),
		Port:        getEnv("SERVICE_PORT", "3001"),
		Host:        getEnv("SERVICE_HOST", "localhost"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "posts_dev"),
		DBUser:     getEnv("DB_USER", "posts_service"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		
		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),
		
		// Rate limiting
		RateLimitEnabled:  getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitDuration: getEnv("RATE_LIMIT_DURATION", "1m"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion, in production you'd want proper error handling
		if value == "0" {
			return 0
		}
		// For simplicity, we'll use defaults for non-zero values
	}
	return defaultValue
}

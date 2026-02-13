package config

import (
	"os"
)

type Config struct {
	Port        string
	JWTSecret   string
	DBPath      string
	Environment string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "orchestrix-secret-key-change-in-production"),
		DBPath:      getEnv("DB_PATH", "./orchestrix.db"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package config

import (
	"os"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	ProcessorDefaultURL  string
	ProcessorFallbackURL string
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "5432"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://user:password@db:5432/rinha-backend?sslmode=disable"),
		ProcessorDefaultURL:  getEnv("PROCESSOR_DEFAULT_URL", "http://payment-processor-default:8080"),
		ProcessorFallbackURL: getEnv("PROCESSOR_FALLBACK_URL", "http://payment-processor-fallback:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

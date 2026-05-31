package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds all environment variables in a strongly-typed struct
type AppConfig struct {
	Port                  string
	DatabaseURL           string
	BaseURL               string
	Environment           string
	ARN                   string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_SES_FROM_EMAIL    string
	AWS_REGION            string
	JWT_SECRET            string
}

// Envs is a global variable holding the loaded configuration
var Envs AppConfig

// InitConfig loads the .env file and populates the AppConfig struct.
// It should be called exactly once at the very start of main().
func InitConfig() {
	// 1. Try to load .env file (Only useful for local dev)
	err := godotenv.Load()
	if err != nil {
		log.Println("[CONFIG] No .env file found. Falling back to system environment variables.")
	}

	// 2. Map OS variables to our struct
	Envs = AppConfig{
		Port:                  getEnv("PORT", "8000"),
		DatabaseURL:           getEnvOrFatal("DATABASE_URL"), // Fail immediately if DB URL is missing
		BaseURL:               getEnv("BASE_URL", "http://localhost:8000"),
		Environment:           getEnv("ENVIRONMENT", "development"),
		ARN:                   getEnv("ARN", ""),
		AWS_ACCESS_KEY_ID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWS_SECRET_ACCESS_KEY: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWS_SES_FROM_EMAIL:    getEnv("AWS_SES_FROM_EMAIL", ""),
		AWS_REGION:            getEnv("AWS_REGION", "ap-south-1"),
		JWT_SECRET:            getEnvOrFatal("JWT_SECRET"),
	}

	log.Println("[CONFIG] Environment variables loaded successfully")
}

// Helper to get an environment variable or return a fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper to require an environment variable. If missing, the app crashes immediately.
func getEnvOrFatal(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatalf("[CONFIG] FATAL: Required environment variable %s is not set", key)
	return ""
}

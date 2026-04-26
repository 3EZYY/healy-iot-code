package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiryHours     int
	CORSAllowedOrigins []string
	WSDevicePath       string
	WSClientPath       string
}

func LoadConfig() *Config {
	// Attempt to load .env file
	_ = godotenv.Load()
	// Parse CORS allowed origins into a slice
	var corsOrigins []string
	corsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	if corsStr != "" {
		for _, o := range strings.Split(corsStr, ",") {
			corsOrigins = append(corsOrigins, strings.TrimSpace(o))
		}
	}

	cfg := &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", "default-unsafe-secret-for-development-only-please-change-it"),
		JWTExpiryHours:     getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		CORSAllowedOrigins: corsOrigins,
		WSDevicePath:       getEnv("WS_DEVICE_PATH", "/ws/device"),
		WSClientPath:       getEnv("WS_CLIENT_PATH", "/ws/telemetry"),
	}

	// Strict validation for critical fields
	if cfg.DatabaseURL == "" {
		log.Println("WARNING: DATABASE_URL environment variable is not set")
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultVal
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("WARNING: Environment variable %s is not a valid integer. Using default: %d\n", key, defaultVal)
		return defaultVal
	}
	return value
}

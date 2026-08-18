package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DatabaseURL      string
	JWTSecret        string
	FrontendURL      string
	TokenExpiryHours int
	UploadDir        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:5173"),
		TokenExpiryHours: getEnvInt("TOKEN_EXPIRY_HOURS", 24),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
	}

	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is required")
		return nil, ErrMissingDatabaseURL
	}
	if cfg.JWTSecret == "" {
		slog.Error("JWT_SECRET is required")
		return nil, ErrMissingJWTSecret
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

var (
	ErrMissingDatabaseURL = &ConfigError{Key: "DATABASE_URL"}
	ErrMissingJWTSecret   = &ConfigError{Key: "JWT_SECRET"}
)

type ConfigError struct {
	Key string
}

func (e *ConfigError) Error() string {
	return "missing required configuration: " + e.Key
}

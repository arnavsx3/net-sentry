package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	GinMode         string
	DatabaseURL     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		ReadTimeout:     getEnvDurationSeconds("HTTP_READ_TIMEOUT_SEC", 10),
		WriteTimeout:    getEnvDurationSeconds("HTTP_WRITE_TIMEOUT_SEC", 10),
		ShutdownTimeout: getEnvDurationSeconds("HTTP_SHUTDOWN_TIMEOUT_SEC", 10),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvDurationSeconds(key string, fallback int) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return time.Duration(fallback) * time.Second
	}

	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return time.Duration(fallback) * time.Second
	}

	return time.Duration(seconds) * time.Second
}
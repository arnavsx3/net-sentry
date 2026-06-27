package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	GinMode         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
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

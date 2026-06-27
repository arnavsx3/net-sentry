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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
	}
}
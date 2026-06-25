package main

import (
    "log"

    "github.com/arnavsx3/net-sentry/backend/internal/config"
    "github.com/arnavsx3/net-sentry/backend/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)

	log.Printf("starting NetSentry backend on %s", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
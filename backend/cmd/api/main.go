package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arnavsx3/net-sentry/backend/internal/config"
	"github.com/arnavsx3/net-sentry/backend/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)

	log.Printf("starting NetSentry backend on port=%s mode=%s", cfg.Port, cfg.GinMode)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}
}

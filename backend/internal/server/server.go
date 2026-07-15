package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/config"
	"github.com/arnavsx3/net-sentry/backend/internal/db"
	"github.com/arnavsx3/net-sentry/backend/internal/handlers"
	"github.com/arnavsx3/net-sentry/backend/internal/repository"
)

type Server struct {
	httpServer *http.Server
	dbClient   *db.Client
}

func New(cfg config.Config, dbClient *db.Client) *Server {
	gin.SetMode(cfg.GinMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	telemetryRepo := repository.NewTelemetryRepository(dbClient)

	engine.GET("/healthz", handlers.HealthCheck)
	engine.GET("/readyz", handlers.ReadinessCheck(dbClient))

	api := engine.Group("/api/v1")
	{
		api.GET("/health", handlers.HealthCheck)
		api.POST("/telemetry", handlers.IngestTelemetry(telemetryRepo))
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		dbClient:   dbClient,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
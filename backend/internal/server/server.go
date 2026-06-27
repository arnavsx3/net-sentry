package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/config"
	"github.com/arnavsx3/net-sentry/backend/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg config.Config) *Server {
	gin.SetMode(cfg.GinMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/healthz", handlers.HealthCheck)
	engine.GET("/readyz", handlers.ReadinessCheck)

	api := engine.Group("/api/v1")
	{
		api.GET("/health", handlers.HealthCheck)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
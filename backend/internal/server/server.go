package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/config"
	"github.com/arnavsx3/net-sentry/backend/internal/handlers"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func New(cfg config.Config) *Server {
	engine := gin.Default()

	engine.GET("/health", handlers.HealthCheck)

	return &Server{
		engine: engine,
		port:   cfg.Port,
	}
}

func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.port))
}

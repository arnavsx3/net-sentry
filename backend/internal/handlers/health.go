package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "net-sentry-backend",
	})
}

func ReadinessCheck(checker ReadinessChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := checker.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"service": "net-sentry-backend",
				"error":   "database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": "net-sentry-backend",
		})
	}
}
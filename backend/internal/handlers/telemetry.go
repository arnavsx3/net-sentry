package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/models"
	"github.com/arnavsx3/net-sentry/backend/internal/repository"
)

func IngestTelemetry(repo *repository.TelemetryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TelemetryIngestRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid telemetry payload",
				"details": err.Error(),
			})
			return
		}

		if err := repo.InsertTelemetry(c.Request.Context(), req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to persist telemetry",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"status":       "accepted",
			"agent_id":     req.AgentID,
			"target_host":  req.Target.Host,
			"trace_hops":   len(req.Trace),
			"probe_status": req.Probe.Status,
		})
	}
}
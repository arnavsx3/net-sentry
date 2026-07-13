package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/models"
)

func IngestTelemetry(c *gin.Context) {
	var req models.TelemetryIngestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid telemetry payload",
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

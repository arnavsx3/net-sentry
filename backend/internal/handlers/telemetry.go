package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/arnavsx3/net-sentry/backend/internal/models"
	"github.com/arnavsx3/net-sentry/backend/internal/realtime"
	"github.com/arnavsx3/net-sentry/backend/internal/repository"
)

type TelemetryBroadcastPayload struct {
	AgentID    string            `json:"agent_id"`
	TargetHost string            `json:"target_host"`
	ObservedAt string            `json:"observed_at"`
	Status     string            `json:"status"`
	LatencyMs  float64           `json:"latency_ms"`
	PacketLoss float64           `json:"packet_loss"`
	Trace      []models.TraceHop `json:"trace"`
}

func IngestTelemetry(repo *repository.TelemetryRepository, hub *realtime.Hub) gin.HandlerFunc {
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

		if hub != nil {
			observedAt := req.Timestamp
			if parsed, err := time.Parse(time.RFC3339, req.Timestamp); err == nil {
				observedAt = parsed.UTC().Format(time.RFC3339)
			}

			hub.Publish("telemetry.received", TelemetryBroadcastPayload{
				AgentID:    req.AgentID,
				TargetHost: req.Target.Host,
				ObservedAt: observedAt,
				Status:     req.Probe.Status,
				LatencyMs:  req.Probe.LatencyMs,
				PacketLoss: req.Probe.PacketLoss,
				Trace:      req.Trace,
			})
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

func GetTargetHistory(repo *repository.TelemetryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Param("host")
		if host == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "target host is required",
			})
			return
		}

		limit := 50
		if raw := c.Query("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid limit",
				})
				return
			}
			limit = parsed
		}

		items, err := repo.GetTargetHistory(c.Request.Context(), host, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to fetch target history",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"target_host": host,
			"count":       len(items),
			"results":     items,
		})
	}
}
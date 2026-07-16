package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/db"
	"github.com/arnavsx3/net-sentry/backend/internal/models"
)

type TelemetryRepository struct {
	orm *gorm.DB
}

type ProbeHistoryItem struct {
	ObservedAt time.Time `json:"observed_at"`
	LatencyMs  float64   `json:"latency_ms"`
	PacketLoss float64   `json:"packet_loss"`
	Status     string    `json:"status"`
}

func NewTelemetryRepository(dbClient *db.Client) *TelemetryRepository {
	return &TelemetryRepository{orm: dbClient.ORM}
}

func (r *TelemetryRepository) InsertTelemetry(ctx context.Context, req models.TelemetryIngestRequest) error {
	observedAt, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return fmt.Errorf("parse timestamp: %w", err)
	}

	return r.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		agent := db.Agent{AgentKey: req.AgentID}
		if err := tx.Where("agent_key = ?", req.AgentID).FirstOrCreate(&agent).Error; err != nil {
			return err
		}

		target := db.Target{Host: req.Target.Host}
		if err := tx.Where("host = ?", req.Target.Host).FirstOrCreate(&target).Error; err != nil {
			return err
		}

		probeResult := db.ProbeResult{
			AgentID:    agent.ID,
			TargetID:   target.ID,
			ObservedAt: observedAt,
			LatencyMs:  req.Probe.LatencyMs,
			PacketLoss: req.Probe.PacketLoss,
			Status:     req.Probe.Status,
		}

		if err := tx.Create(&probeResult).Error; err != nil {
			return err
		}

		for _, hop := range req.Trace {
			traceHop := db.TracerouteHop{
				ProbeResultID: probeResult.ID,
				HopNumber:     hop.Hop,
				Address:       hop.Address,
				RTTMs:         hop.RTTMs,
			}

			if err := tx.Create(&traceHop).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TelemetryRepository) GetTargetHistory(ctx context.Context, host string, limit int) ([]ProbeHistoryItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}

	var target db.Target
	if err := r.orm.WithContext(ctx).
		Where("host = ?", host).
		First(&target).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []ProbeHistoryItem{}, nil
		}
		return nil, err
	}

	var probeResults []db.ProbeResult
	if err := r.orm.WithContext(ctx).
		Where("target_id = ?", target.ID).
		Order("observed_at desc").
		Limit(limit).
		Find(&probeResults).Error; err != nil {
		return nil, err
	}

	items := make([]ProbeHistoryItem, 0, len(probeResults))
	for _, pr := range probeResults {
		items = append(items, ProbeHistoryItem{
			ObservedAt: pr.ObservedAt,
			LatencyMs:  pr.LatencyMs,
			PacketLoss: pr.PacketLoss,
			Status:     pr.Status,
		})
	}

	return items, nil
}
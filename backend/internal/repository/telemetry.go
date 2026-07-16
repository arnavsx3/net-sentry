package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/config"
	"github.com/arnavsx3/net-sentry/backend/internal/db"
	"github.com/arnavsx3/net-sentry/backend/internal/models"
)

const (
	alertTypeLatency    = "latency"
	alertTypePacketLoss = "packet_loss"
	alertTypeTargetDown = "target_down"

	alertSeverityWarning  = "warning"
	alertSeverityCritical = "critical"
)

type TelemetryRepository struct {
	orm                 *gorm.DB
	latencyThresholdMs  float64
	packetLossThreshold float64
}

type ProbeHistoryItem struct {
	ObservedAt time.Time `json:"observed_at"`
	LatencyMs  float64   `json:"latency_ms"`
	PacketLoss float64   `json:"packet_loss"`
	Status     string    `json:"status"`
}

func NewTelemetryRepository(cfg config.Config, dbClient *db.Client) *TelemetryRepository {
	return &TelemetryRepository{
		orm:                 dbClient.ORM,
		latencyThresholdMs:  cfg.AlertLatencyMsThreshold,
		packetLossThreshold: cfg.AlertPacketLossThreshold,
	}
}

func (r *TelemetryRepository) InsertTelemetry(ctx context.Context, req models.TelemetryIngestRequest) error {
	observedAt, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return fmt.Errorf("parse timestamp: %w", err)
	}

	return r.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		agent := db.Agent{AgentKey: req.AgentID}
		if err := tx.Where(&db.Agent{AgentKey: req.AgentID}).FirstOrCreate(&agent).Error; err != nil {
			return err
		}

		target := db.Target{Host: req.Target.Host}
		if err := tx.Where(&db.Target{Host: req.Target.Host}).FirstOrCreate(&target).Error; err != nil {
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

		if err := r.reconcileAlerts(tx, target, probeResult); err != nil {
			return err
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
		Where(&db.Target{Host: host}).
		First(&target).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []ProbeHistoryItem{}, nil
		}
		return nil, err
	}

	var probeResults []db.ProbeResult
	if err := r.orm.WithContext(ctx).
		Where(&db.ProbeResult{TargetID: target.ID}).
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

func (r *TelemetryRepository) reconcileAlerts(tx *gorm.DB, target db.Target, probeResult db.ProbeResult) error {
	if probeResult.Status == "down" {
		if err := r.upsertActiveAlert(
			tx,
			target.ID,
			probeResult,
			alertTypeTargetDown,
			alertSeverityCritical,
			fmt.Sprintf("target %s is down", target.Host),
		); err != nil {
			return err
		}

		if err := r.resolveAlertType(tx, target.ID, alertTypeLatency, probeResult.ObservedAt); err != nil {
			return err
		}

		if err := r.resolveAlertType(tx, target.ID, alertTypePacketLoss, probeResult.ObservedAt); err != nil {
			return err
		}

		return nil
	}

	if err := r.resolveAlertType(tx, target.ID, alertTypeTargetDown, probeResult.ObservedAt); err != nil {
		return err
	}

	if probeResult.LatencyMs >= r.latencyThresholdMs {
		if err := r.upsertActiveAlert(
			tx,
			target.ID,
			probeResult,
			alertTypeLatency,
			alertSeverityWarning,
			fmt.Sprintf("latency %.2fms exceeded threshold %.2fms for %s", probeResult.LatencyMs, r.latencyThresholdMs, target.Host),
		); err != nil {
			return err
		}
	} else {
		if err := r.resolveAlertType(tx, target.ID, alertTypeLatency, probeResult.ObservedAt); err != nil {
			return err
		}
	}

	if probeResult.PacketLoss >= r.packetLossThreshold {
		if err := r.upsertActiveAlert(
			tx,
			target.ID,
			probeResult,
			alertTypePacketLoss,
			alertSeverityCritical,
			fmt.Sprintf("packet loss %.2f%% exceeded threshold %.2f%% for %s", probeResult.PacketLoss, r.packetLossThreshold, target.Host),
		); err != nil {
			return err
		}
	} else {
		if err := r.resolveAlertType(tx, target.ID, alertTypePacketLoss, probeResult.ObservedAt); err != nil {
			return err
		}
	}

	return nil
}

func (r *TelemetryRepository) upsertActiveAlert(
	tx *gorm.DB,
	targetID uint,
	probeResult db.ProbeResult,
	alertType string,
	severity string,
	message string,
) error {
	var alert db.Alert

	err := tx.
		Where(map[string]any{
			"target_id":   targetID,
			"type":        alertType,
			"resolved_at": nil,
		}).
		First(&alert).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			alert = db.Alert{
				TargetID:      targetID,
				ProbeResultID: probeResult.ID,
				Type:          alertType,
				Severity:      severity,
				Message:       message,
				TriggeredAt:   probeResult.ObservedAt,
			}

			return tx.Create(&alert).Error
		}

		return err
	}

	return tx.Model(&alert).Updates(map[string]any{
		"probe_result_id": probeResult.ID,
		"severity":        severity,
		"message":         message,
	}).Error
}

func (r *TelemetryRepository) resolveAlertType(tx *gorm.DB, targetID uint, alertType string, resolvedAt time.Time) error {
	return tx.Model(&db.Alert{}).
		Where(map[string]any{
			"target_id":   targetID,
			"type":        alertType,
			"resolved_at": nil,
		}).
		Updates(map[string]any{
			"resolved_at": resolvedAt,
		}).Error
}
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/db"
)

type AlertRepository struct {
	orm *gorm.DB
}

type CurrentAlertItem struct {
	ID            uint       `json:"id"`
	TargetHost    string     `json:"target_host"`
	Type          string     `json:"type"`
	Severity      string     `json:"severity"`
	Message       string     `json:"message"`
	TriggeredAt   time.Time  `json:"triggered_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
	ObservedAt    time.Time  `json:"observed_at"`
	LatencyMs     float64    `json:"latency_ms"`
	PacketLoss    float64    `json:"packet_loss"`
	ProbeStatus   string     `json:"probe_status"`
	ProbeResultID uint       `json:"probe_result_id"`
}

func NewAlertRepository(dbClient *db.Client) *AlertRepository {
	return &AlertRepository{orm: dbClient.ORM}
}

func (r *AlertRepository) GetCurrentAlerts(ctx context.Context, limit int) ([]CurrentAlertItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var alerts []db.Alert
	if err := r.orm.WithContext(ctx).
		Preload("Target").
		Preload("ProbeResult").
		Where(map[string]any{
			"resolved_at": nil,
		}).
		Order("triggered_at desc").
		Limit(limit).
		Find(&alerts).Error; err != nil {
		return nil, err
	}

	items := make([]CurrentAlertItem, 0, len(alerts))
	for _, alert := range alerts {
		items = append(items, CurrentAlertItem{
			ID:            alert.ID,
			TargetHost:    alert.Target.Host,
			Type:          alert.Type,
			Severity:      alert.Severity,
			Message:       alert.Message,
			TriggeredAt:   alert.TriggeredAt,
			ResolvedAt:    alert.ResolvedAt,
			ObservedAt:    alert.ProbeResult.ObservedAt,
			LatencyMs:     alert.ProbeResult.LatencyMs,
			PacketLoss:    alert.ProbeResult.PacketLoss,
			ProbeStatus:   alert.ProbeResult.Status,
			ProbeResultID: alert.ProbeResultID,
		})
	}

	return items, nil
}
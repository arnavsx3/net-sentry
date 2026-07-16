package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/db"
)

type TargetRepository struct {
	orm *gorm.DB
}

type TargetCurrentAlert struct {
	Type       string `json:"type"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
}

type TargetCurrentItem struct {
	TargetHost        string               `json:"target_host"`
	Status            string               `json:"status"`
	LatencyMs         float64              `json:"latency_ms"`
	PacketLoss        float64              `json:"packet_loss"`
	ObservedAt        *time.Time           `json:"observed_at"`
	ActiveAlertCount  int                  `json:"active_alert_count"`
	ActiveAlerts      []TargetCurrentAlert `json:"active_alerts"`
}

func NewTargetRepository(dbClient *db.Client) *TargetRepository {
	return &TargetRepository{orm: dbClient.ORM}
}

func (r *TargetRepository) GetCurrentTargets(ctx context.Context, limit int) ([]TargetCurrentItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var targets []db.Target
	if err := r.orm.WithContext(ctx).
		Order("host asc").
		Limit(limit).
		Find(&targets).Error; err != nil {
		return nil, err
	}

	items := make([]TargetCurrentItem, 0, len(targets))

	for _, target := range targets {
		item := TargetCurrentItem{
			TargetHost:   target.Host,
			Status:       "unknown",
			ActiveAlerts: []TargetCurrentAlert{},
		}

		var latestProbe db.ProbeResult
		err := r.orm.WithContext(ctx).
			Where(&db.ProbeResult{TargetID: target.ID}).
			Order("observed_at desc").
			First(&latestProbe).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if err == nil {
			observedAt := latestProbe.ObservedAt
			item.Status = latestProbe.Status
			item.LatencyMs = latestProbe.LatencyMs
			item.PacketLoss = latestProbe.PacketLoss
			item.ObservedAt = &observedAt
		}

		var activeAlerts []db.Alert
		if err := r.orm.WithContext(ctx).
			Where(map[string]any{
				"target_id":   target.ID,
				"resolved_at": nil,
			}).
			Order("triggered_at desc").
			Find(&activeAlerts).Error; err != nil {
			return nil, err
		}

		item.ActiveAlertCount = len(activeAlerts)
		for _, alert := range activeAlerts {
			item.ActiveAlerts = append(item.ActiveAlerts, TargetCurrentAlert{
				Type:        alert.Type,
				Severity:    alert.Severity,
				Message:     alert.Message,
				TriggeredAt: alert.TriggeredAt,
			})
		}

		items = append(items, item)
	}

	return items, nil
}
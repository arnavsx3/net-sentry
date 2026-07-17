package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/db"
)

type TracerouteRepository struct {
	orm *gorm.DB
}

type TracerouteHopItem struct {
	Hop     int     `json:"hop"`
	Address string  `json:"address"`
	RTTMs   float64 `json:"rtt_ms"`
}

type LatestTracerouteItem struct {
	TargetHost   string               `json:"target_host"`
	ObservedAt   *time.Time           `json:"observed_at"`
	ProbeStatus  string               `json:"probe_status"`
	LatencyMs    float64              `json:"latency_ms"`
	PacketLoss   float64              `json:"packet_loss"`
	Hops         []TracerouteHopItem  `json:"hops"`
}

func NewTracerouteRepository(dbClient *db.Client) *TracerouteRepository {
	return &TracerouteRepository{orm: dbClient.ORM}
}

func (r *TracerouteRepository) GetLatestTraceroute(ctx context.Context, host string) (LatestTracerouteItem, error) {
	result := LatestTracerouteItem{
		TargetHost: host,
		Hops:       []TracerouteHopItem{},
	}

	var target db.Target
	if err := r.orm.WithContext(ctx).
		Where(&db.Target{Host: host}).
		First(&target).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, nil
		}
		return result, err
	}

	var probeResult db.ProbeResult
	if err := r.orm.WithContext(ctx).
		Where(&db.ProbeResult{TargetID: target.ID}).
		Order("observed_at desc").
		First(&probeResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, nil
		}
		return result, err
	}

	observedAt := probeResult.ObservedAt
	result.ObservedAt = &observedAt
	result.ProbeStatus = probeResult.Status
	result.LatencyMs = probeResult.LatencyMs
	result.PacketLoss = probeResult.PacketLoss

	var hops []db.TracerouteHop
	if err := r.orm.WithContext(ctx).
		Where(&db.TracerouteHop{ProbeResultID: probeResult.ID}).
		Order("hop_number asc").
		Find(&hops).Error; err != nil {
		return result, err
	}

	result.Hops = make([]TracerouteHopItem, 0, len(hops))
	for _, hop := range hops {
		result.Hops = append(result.Hops, TracerouteHopItem{
			Hop:     hop.HopNumber,
			Address: hop.Address,
			RTTMs:   hop.RTTMs,
		})
	}

	return result, nil
}
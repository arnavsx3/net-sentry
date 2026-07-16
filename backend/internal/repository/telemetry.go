type ProbeHistoryItem struct {
	ObservedAt time.Time `json:"observed_at"`
	LatencyMs  float64   `json:"latency_ms"`
	PacketLoss float64   `json:"packet_loss"`
	Status     string    `json:"status"`
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
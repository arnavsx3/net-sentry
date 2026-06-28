package models

type TargetPayload struct {
	Host string `json:"host" binding:"required"`
}

type ProbePayload struct {
	LatencyMs  float64 `json:"latency_ms" binding:"gte=0"`
	PacketLoss float64 `json:"packet_loss" binding:"gte=0,lte=100"`
	Status     string  `json:"status" binding:"required,oneof=healthy degraded down"`
}
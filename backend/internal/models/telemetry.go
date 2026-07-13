package models

type TelemetryIngestRequest struct {
	AgentID   string        `json:"agent_id" binding:"required"`
	Timestamp string        `json:"timestamp" binding:"required"`
	Target    TargetPayload `json:"target" binding:"required"`
	Probe     ProbePayload  `json:"probe" binding:"required"`
	Trace     []TraceHop    `json:"trace"`
}

type TargetPayload struct {
	Host string `json:"host" binding:"required"`
}

type ProbePayload struct {
	LatencyMs  float64 `json:"latency_ms" binding:"gte=0"`
	PacketLoss float64 `json:"packet_loss" binding:"gte=0,lte=100"`
	Status     string  `json:"status" binding:"required,oneof=healthy degraded down"`
}

type TraceHop struct {
	Hop     int     `json:"hop" binding:"required,gte=1"`
	Address string  `json:"address" binding:"required"`
	RTTMs   float64 `json:"rtt_ms" binding:"gte=0"`
}
package db

import "time"

type Agent struct {
	ID        uint      `gorm:"primaryKey"`
	AgentKey  string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

type Target struct {
	ID        uint      `gorm:"primaryKey"`
	Host      string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

type ProbeResult struct {
	ID         uint      `gorm:"primaryKey"`
	AgentID    uint      `gorm:"not null;index"`
	TargetID   uint      `gorm:"not null;index"`
	ObservedAt time.Time `gorm:"not null;index"`
	LatencyMs  float64   `gorm:"not null"`
	PacketLoss float64   `gorm:"not null"`
	Status     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
}

type TracerouteHop struct {
	ID            uint      `gorm:"primaryKey"`
	ProbeResultID uint      `gorm:"not null;index"`
	HopNumber     int       `gorm:"not null"`
	Address       string    `gorm:"not null"`
	RTTMs         float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

type Alert struct {
	ID            uint        `gorm:"primaryKey"`
	TargetID      uint        `gorm:"not null;index"`
	ProbeResultID uint        `gorm:"not null;index"`
	Type          string      `gorm:"not null;index"`
	Severity      string      `gorm:"not null"`
	Message       string      `gorm:"not null"`
	TriggeredAt   time.Time   `gorm:"not null;index"`
	ResolvedAt    *time.Time  `gorm:"index"`
	CreatedAt     time.Time   `gorm:"not null"`
	UpdatedAt     time.Time   `gorm:"not null"`
	Target        Target      `gorm:"foreignKey:TargetID"`
	ProbeResult   ProbeResult `gorm:"foreignKey:ProbeResultID"`
}
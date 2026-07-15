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
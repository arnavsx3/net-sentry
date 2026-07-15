package db

import "time"

type Agent struct {
	ID        uint      `gorm:"primaryKey"`
	AgentKey  string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

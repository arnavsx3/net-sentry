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

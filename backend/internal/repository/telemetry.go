package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/arnavsx3/net-sentry/backend/internal/db"
	"github.com/arnavsx3/net-sentry/backend/internal/models"
)

type TelemetryRepository struct {
	orm *gorm.DB
}
package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	ORM   *gorm.DB
	sqlDB *sql.DB
}

func New(ctx context.Context, databaseURL string) (*Client, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	orm, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, err
	}

	if err := orm.AutoMigrate(
		&Agent{},
		&Target{},
		&ProbeResult{},
		&TracerouteHop{},
		&Alert{},
	); err != nil {
		return nil, err
	}

	return &Client{
		ORM:   orm,
		sqlDB: sqlDB,
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.sqlDB == nil {
		return nil
	}

	return c.sqlDB.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.sqlDB.PingContext(ctx)
}
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool membuat connection pool ke TimescaleDB/PostgreSQL menggunakan pgxpool.
// Pool config mengikuti best practices untuk high-frequency IoT data ingestion:
//   - MaxConns: 25 (cukup untuk concurrent WS writes + REST reads)
//   - MinConns: 5 (keep warm connections untuk low-latency first queries)
//   - MaxConnLifetime: 1 jam (cegah stale connections)
//   - MaxConnIdleTime: 30 menit
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verifikasi koneksi berhasil
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

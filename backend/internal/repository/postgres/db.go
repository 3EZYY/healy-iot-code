package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB membungkus pgxpool.Pool untuk koneksi ke Supabase (PostgreSQL/TimescaleDB).
// Koneksi menggunakan DATABASE_URL dari environment variable — tidak ada hardcoded credentials.
//
// Referensi: HEALY_Master_Blueprint.md Section 7.2
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB membuat koneksi pool ke Supabase menggunakan DATABASE_URL dari env.
//
// Format Supabase connection string:
//
//	postgresql://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres
//
// Pool config dioptimalkan untuk high-frequency IoT telemetry ingestion:
//   - MaxConns: 25 — cukup untuk concurrent WS writes + REST reads
//   - MinConns: 5  — warm connections untuk low-latency first queries
//   - MaxConnLifetime: 1 jam — cegah stale connections ke Supabase pooler
//   - MaxConnIdleTime: 30 menit
func NewDB(ctx context.Context) (*DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	// Pool settings — optimized for IoT sensor data volume
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verifikasi koneksi ke Supabase berhasil
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping Supabase database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close menutup semua koneksi di pool. Panggil saat graceful shutdown.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Health memeriksa apakah koneksi ke database masih aktif.
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

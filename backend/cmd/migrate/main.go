// cmd/migrate/main.go
//
// HEALY Database Migration Runner
// ================================
// Membaca DATABASE_URL dari .env, connect ke Supabase via pgx/v5,
// lalu mengeksekusi file SQL migration secara berurutan.
//
// Usage:
//
//	# Dari direktori backend/
//	go run ./cmd/migrate/main.go
//
// Referensi: HEALY_Master_Blueprint.md Section 7.2
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// migrationFiles adalah daftar file SQL yang akan dieksekusi secara berurutan.
// Tambahkan file baru di sini saat ada migration berikutnya.
var migrationFiles = []string{
	"migrations/000001_init_healy.up.sql",
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	log.Info("HEALY Migration Runner starting")

	// ── 1. Resolve project root ────────────────────────────────────
	// Migration runner dijalankan dari direktori backend/.
	// Kita resolve path relatif terhadap lokasi file ini agar
	// tidak bergantung pada working directory saat go run dipanggil.
	projectRoot := resolveProjectRoot()
	log.Info("project root resolved", slog.String("path", projectRoot))

	// ── 2. Load .env ───────────────────────────────────────────────
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		// .env tidak wajib ada jika DATABASE_URL sudah di-set di environment
		log.Warn(".env file not found, falling back to system environment",
			slog.String("path", envPath),
			slog.String("error", err.Error()),
		)
	} else {
		log.Info(".env loaded", slog.String("path", envPath))
	}

	// ── 3. Baca DATABASE_URL ───────────────────────────────────────
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Error("DATABASE_URL is not set. Set it in backend/.env or as an environment variable.")
		os.Exit(1)
	}
	log.Info("DATABASE_URL found", slog.String("hint", maskDSN(databaseURL)))

	// ── 4. Connect ke Supabase ─────────────────────────────────────
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Error("failed to connect to Supabase",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Verifikasi koneksi aktif
	if err := conn.Ping(ctx); err != nil {
		log.Error("database ping failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("connected to Supabase successfully")

	// ── 5. Eksekusi setiap migration file ─────────────────────────
	successCount := 0
	for _, relPath := range migrationFiles {
		sqlPath := filepath.Join(projectRoot, relPath)

		log.Info("running migration", slog.String("file", relPath))

		if err := runMigration(ctx, conn, sqlPath, log); err != nil {
			log.Error("migration failed",
				slog.String("file", relPath),
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}

		log.Info("migration completed", slog.String("file", relPath))
		successCount++
	}

	// ── 6. Ringkasan ───────────────────────────────────────────────
	log.Info("all migrations completed successfully",
		slog.Int("total", successCount),
	)
}

// runMigration membaca file SQL dan mengeksekusinya dalam satu transaksi.
// Jika ada error di tengah eksekusi, seluruh transaksi di-rollback.
func runMigration(ctx context.Context, conn *pgx.Conn, sqlPath string, log *slog.Logger) error {
	// Baca file SQL
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %q: %w", sqlPath, err)
	}

	sqlContent := string(sqlBytes)
	if len(sqlContent) == 0 {
		return fmt.Errorf("SQL file is empty: %q", sqlPath)
	}

	log.Debug("SQL file read",
		slog.String("path", sqlPath),
		slog.Int("bytes", len(sqlBytes)),
	)

	// Eksekusi dalam transaksi agar atomic
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Pastikan rollback dipanggil jika terjadi panic atau error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // re-panic setelah rollback
		}
	}()

	if _, err := tx.Exec(ctx, sqlContent); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// resolveProjectRoot mengembalikan path absolut ke direktori backend/.
// Menggunakan runtime.Caller untuk mendapatkan lokasi file ini,
// sehingga tidak bergantung pada working directory saat dijalankan.
func resolveProjectRoot() string {
	// Dapatkan path file ini saat compile time
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		// Fallback: gunakan working directory
		wd, _ := os.Getwd()
		return wd
	}

	// filename = .../backend/cmd/migrate/main.go
	// Naik 3 level: migrate/ → cmd/ → backend/
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

// maskDSN menyembunyikan password dari connection string untuk logging.
// Contoh: "postgresql://user:SECRET@host/db" → "postgresql://user:***@host/db"
func maskDSN(dsn string) string {
	// Cari posisi "://" dan "@" untuk isolasi credentials
	schemeEnd := -1
	for i := 0; i < len(dsn)-2; i++ {
		if dsn[i] == ':' && dsn[i+1] == '/' && dsn[i+2] == '/' {
			schemeEnd = i + 3
			break
		}
	}
	if schemeEnd == -1 {
		return "[masked DSN]"
	}

	atPos := -1
	for i := schemeEnd; i < len(dsn); i++ {
		if dsn[i] == '@' {
			atPos = i
			break
		}
	}
	if atPos == -1 {
		return dsn[:schemeEnd] + "***"
	}

	// Cari ":" antara schemeEnd dan atPos (pemisah user:password)
	colonPos := -1
	for i := schemeEnd; i < atPos; i++ {
		if dsn[i] == ':' {
			colonPos = i
			break
		}
	}
	if colonPos == -1 {
		return dsn[:schemeEnd] + "***" + dsn[atPos:]
	}

	return dsn[:colonPos+1] + "***" + dsn[atPos:]
}

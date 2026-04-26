-- HEALY Backend — Initial Schema Migration
-- Database: TimescaleDB (PostgreSQL extension untuk time-series)
-- Referensi: HEALY_Master_Blueprint.md Section 6
--
-- Jalankan migration ini setelah docker-compose up:
--   psql -h localhost -U healy_user -d healy_db -f migrations/001_init_schema.up.sql

-- ============================================================
-- 1. Aktifkan TimescaleDB extension
-- ============================================================
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================================
-- 2. Tabel users — untuk login dashboard
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    password   VARCHAR(255) NOT NULL,   -- bcrypt hash, cost >= 12
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- 3. Tabel telemetry_records — hypertable time-series
-- ============================================================
-- Composite primary key (id, recorded_at) diperlukan oleh TimescaleDB
-- agar partitioning berdasarkan waktu bisa bekerja.
CREATE TABLE IF NOT EXISTS telemetry_records (
    id             UUID DEFAULT gen_random_uuid(),
    device_id      VARCHAR(50) NOT NULL,
    recorded_at    TIMESTAMPTZ NOT NULL,        -- Timestamp dari ESP32
    temperature    DECIMAL(4,1) NOT NULL,       -- Celsius, presisi 1 desimal
    bpm            SMALLINT NOT NULL,           -- Beats per minute
    spo2           SMALLINT NOT NULL,           -- Persentase SpO2
    temp_status    VARCHAR(10) NOT NULL,        -- NORMAL / WARNING / CRITICAL
    spo2_status    VARCHAR(10) NOT NULL,
    overall_status VARCHAR(10) NOT NULL,
    PRIMARY KEY (id, recorded_at)
);

-- Konversi ke hypertable dengan chunk interval 1 hari.
-- Optimal untuk volume ~50K-100K rows/hari per device.
SELECT create_hypertable('telemetry_records', 'recorded_at',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Index untuk query history by device + time range (paling sering dipakai)
CREATE INDEX IF NOT EXISTS idx_telemetry_device_time
    ON telemetry_records (device_id, recorded_at DESC);

-- ============================================================
-- 4. Tabel alert_logs — log setiap alert yang dipicu
-- ============================================================
CREATE TABLE IF NOT EXISTS alert_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id    VARCHAR(50) NOT NULL,
    alert_type   VARCHAR(20) NOT NULL,     -- TEMP_CRITICAL, SPO2_WARNING, dll
    value        DECIMAL(5,2) NOT NULL,    -- Nilai sensor saat alert dipicu
    status       VARCHAR(10) NOT NULL,     -- WARNING / CRITICAL
    triggered_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index untuk query alert deduplication (cek alert terakhir per device)
CREATE INDEX IF NOT EXISTS idx_alert_device_time
    ON alert_logs (device_id, triggered_at DESC);

-- ============================================================
-- 5. Tabel device_settings — konfigurasi threshold per device
-- ============================================================
CREATE TABLE IF NOT EXISTS device_settings (
    device_id     VARCHAR(50) PRIMARY KEY,
    temp_warn_max DECIMAL(4,1) DEFAULT 37.5,   -- Batas atas suhu WARNING
    temp_crit_max DECIMAL(4,1) DEFAULT 38.5,   -- Batas atas suhu CRITICAL
    spo2_warn_min SMALLINT DEFAULT 94,          -- Batas bawah SpO2 WARNING
    spo2_crit_min SMALLINT DEFAULT 90,          -- Batas bawah SpO2 CRITICAL
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- 6. Seed data — default user untuk development
-- ============================================================
-- Password: healy123 (bcrypt hash, cost 12)
-- Generate ulang dengan: htpasswd -nbBC 12 "" healy123 | cut -d: -f2
INSERT INTO users (username, password)
VALUES ('admin', '$2a$12$LJ3m4ys.aGW0ZFEqLEwRhOaCDRVf6bDE1GGOmLGjn0cR2heEz1OHi')
ON CONFLICT (username) DO NOTHING;

-- Default threshold settings untuk device healy-001
INSERT INTO device_settings (device_id)
VALUES ('healy-001')
ON CONFLICT (device_id) DO NOTHING;

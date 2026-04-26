-- ============================================================
-- HEALY Backend — Initial Schema Migration
-- Database: Supabase (PostgreSQL) + pg_partman
-- Referensi: HEALY_Master_Blueprint.md Section 6
--
-- CATATAN ARSITEKTUR:
--   Supabase tidak mendukung TimescaleDB extension.
--   Sebagai gantinya, kita menggunakan pg_partman untuk
--   declarative time-based table partitioning (daily).
--   pg_partman sudah tersedia di Supabase secara default.
-- ============================================================

-- ============================================================
-- 1. Aktifkan pg_partman extension
--    (sudah tersedia di Supabase, tidak perlu install manual)
--    Note: Supabase menyimpan extension di schema 'extensions'
-- ============================================================
CREATE EXTENSION IF NOT EXISTS pg_partman SCHEMA extensions;

-- ============================================================
-- 2. Tabel users — untuk login dashboard
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    password   VARCHAR(255) NOT NULL,   -- bcrypt hash
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- 3. Tabel telemetry_records — partitioned by recorded_at (daily)
--
--    Menggunakan PARTITION BY RANGE (recorded_at) agar pg_partman
--    bisa mengelola child partitions secara otomatis.
--
--    PENTING: Composite PK harus menyertakan partition key (recorded_at)
--    agar PostgreSQL native partitioning bisa bekerja.
-- ============================================================
CREATE TABLE IF NOT EXISTS telemetry_records (
    id             UUID            NOT NULL DEFAULT gen_random_uuid(),
    device_id      VARCHAR(50)     NOT NULL,
    recorded_at    TIMESTAMPTZ     NOT NULL,   -- dari ESP32 timestamp
    temperature    DECIMAL(4,1)    NOT NULL,
    bpm            SMALLINT        NOT NULL,
    spo2           SMALLINT        NOT NULL,
    temp_status    VARCHAR(10)     NOT NULL,   -- NORMAL/WARNING/CRITICAL
    spo2_status    VARCHAR(10)     NOT NULL,
    overall_status VARCHAR(10)     NOT NULL,
    PRIMARY KEY (id, recorded_at)             -- composite PK wajib untuk partitioning
) PARTITION BY RANGE (recorded_at);

-- ============================================================
-- 4. Daftarkan tabel ke pg_partman untuk daily auto-partitioning
--
--    p_parent_table  : nama tabel parent (schema-qualified)
--    p_control       : kolom yang digunakan sebagai partition key
--    p_type          : 'native' untuk menggunakan native partitioning
--    p_interval      : interval partisi — 'daily' = 1 hari per partisi
--    p_premake       : buat 4 partisi ke depan secara otomatis
-- ============================================================
SELECT extensions.create_parent(
    p_parent_table  => 'public.telemetry_records',
    p_control       => 'recorded_at',
    p_type          => 'range',
    p_interval      => '1 day',
    p_premake       => 4
);

-- ============================================================
-- 5. Index untuk query history by device (pada parent table)
--    pg_partman akan mewarisi index ini ke child partitions.
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_telemetry_device_time
    ON telemetry_records (device_id, recorded_at DESC);

-- ============================================================
-- 6. Tabel alert_logs — log setiap alert yang dipicu
-- ============================================================
CREATE TABLE IF NOT EXISTS alert_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id    VARCHAR(50)  NOT NULL,
    alert_type   VARCHAR(20)  NOT NULL,   -- TEMP_CRITICAL, SPO2_WARNING, dsb
    value        DECIMAL(5,2) NOT NULL,
    status       VARCHAR(10)  NOT NULL,
    triggered_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alert_logs_device_time
    ON alert_logs (device_id, triggered_at DESC);

-- ============================================================
-- 7. Tabel device_settings — konfigurasi threshold per device
-- ============================================================
CREATE TABLE IF NOT EXISTS device_settings (
    device_id     VARCHAR(50)  PRIMARY KEY,
    temp_warn_max DECIMAL(4,1) DEFAULT 37.5,
    temp_crit_max DECIMAL(4,1) DEFAULT 38.5,
    spo2_warn_min SMALLINT     DEFAULT 94,
    spo2_crit_min SMALLINT     DEFAULT 90,
    updated_at    TIMESTAMPTZ  DEFAULT NOW()
);

-- ============================================================
-- 8. Seed data: default device settings untuk healy-001
-- ============================================================
INSERT INTO device_settings (device_id, temp_warn_max, temp_crit_max, spo2_warn_min, spo2_crit_min)
VALUES ('healy-001', 37.5, 38.5, 94, 90)
ON CONFLICT (device_id) DO NOTHING;

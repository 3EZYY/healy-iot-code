-- HEALY Backend — Rollback Initial Schema
-- Drops semua tabel dalam urutan terbalik (foreign key safe).

DROP TABLE IF EXISTS device_settings;
DROP TABLE IF EXISTS alert_logs;
DROP TABLE IF EXISTS telemetry_records;
DROP TABLE IF EXISTS users;

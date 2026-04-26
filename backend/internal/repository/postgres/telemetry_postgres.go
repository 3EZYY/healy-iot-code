package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/repository/interfaces"
)

type telemetryRepository struct {
	pool *pgxpool.Pool
}

// NewTelemetryRepository creates a new instance of TelemetryRepository
// utilizing the pgx/v5 connection pool.
func NewTelemetryRepository(pool *pgxpool.Pool) interfaces.TelemetryRepository {
	return &telemetryRepository{
		pool: pool,
	}
}

// Save inserts a new telemetry record into the telemetry.telemetry_records table.
// CRITICAL: Must use schema-qualified name 'telemetry.telemetry_records'.
func (r *telemetryRepository) Save(ctx context.Context, record domain.TelemetryRecord) error {
	const insertQuery = `
		INSERT INTO telemetry.telemetry_records 
		(device_id, recorded_at, temperature, bpm, spo2, temp_status, spo2_status, overall_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, insertQuery,
		record.DeviceID,
		record.Timestamp,
		record.Sensor.Temperature,
		record.Sensor.BPM,
		record.Sensor.SpO2,
		string(record.Status.Temperature),
		string(record.Status.SpO2),
		string(record.Status.Overall),
	)
	if err != nil {
		return fmt.Errorf("failed to insert telemetry record: %w", err)
	}
	return nil
}

// SaveAlert inserts a critical alert log into the public.alert_logs table.
func (r *telemetryRepository) SaveAlert(ctx context.Context, alert domain.Alert) error {
	const insertAlert = `
		INSERT INTO public.alert_logs 
		(device_id, alert_type, value, status, triggered_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, insertAlert,
		alert.DeviceID,
		string(alert.AlertType),
		alert.Value,
		string(alert.Status),
		alert.TriggeredAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert alert log: %w", err)
	}
	return nil
}

// GetHistory retrieves historical telemetry data aggregated by minute.
func (r *telemetryRepository) GetHistory(ctx context.Context, deviceID string, hours int) ([]domain.TelemetryRecord, error) {
	// Construct the interval string securely (e.g., "24 hours")
	interval := fmt.Sprintf("%d hours", hours)

	const queryHistory = `
		SELECT 
		  DATE_TRUNC('minute', recorded_at) AS bucket,
		  AVG(temperature)::DECIMAL(4,1) AS avg_temp,
		  AVG(bpm)::INT AS avg_bpm,
		  AVG(spo2)::INT AS avg_spo2,
		  MAX(overall_status) AS status
		FROM telemetry.telemetry_records
		WHERE device_id = $1
		  AND recorded_at >= NOW() - $2::interval
		GROUP BY bucket
		ORDER BY bucket DESC
		LIMIT 500
	`

	rows, err := r.pool.Query(ctx, queryHistory, deviceID, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var records []domain.TelemetryRecord
	for rows.Next() {
		var record domain.TelemetryRecord
		var overallStatus string

		// Note: History aggregates sensor data by minute, so we construct a record holding averages
		err := rows.Scan(
			&record.Timestamp,
			&record.Sensor.Temperature,
			&record.Sensor.BPM,
			&record.Sensor.SpO2,
			&overallStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}
		
		record.DeviceID = deviceID
		record.Status.Overall = domain.SensorStatus(overallStatus)
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("history rows error: %w", err)
	}

	return records, nil
}

// GetLatest retrieves the most recent telemetry record for a specific device.
func (r *telemetryRepository) GetLatest(ctx context.Context, deviceID string) (domain.TelemetryRecord, error) {
	const queryLatest = `
		SELECT DISTINCT ON (device_id)
		  device_id, recorded_at, temperature, bpm, spo2,
		  temp_status, spo2_status, overall_status
		FROM telemetry.telemetry_records
		WHERE device_id = $1
		ORDER BY device_id, recorded_at DESC
	`

	var record domain.TelemetryRecord
	var tempStatus, spo2Status, overallStatus string

	err := r.pool.QueryRow(ctx, queryLatest, deviceID).Scan(
		&record.DeviceID,
		&record.Timestamp,
		&record.Sensor.Temperature,
		&record.Sensor.BPM,
		&record.Sensor.SpO2,
		&tempStatus,
		&spo2Status,
		&overallStatus,
	)
	if err != nil {
		return domain.TelemetryRecord{}, fmt.Errorf("failed to query latest record: %w", err)
	}

	record.Status.Temperature = domain.SensorStatus(tempStatus)
	record.Status.SpO2 = domain.SensorStatus(spo2Status)
	record.Status.Overall = domain.SensorStatus(overallStatus)

	return record, nil
}

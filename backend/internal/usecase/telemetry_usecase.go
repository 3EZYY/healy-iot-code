package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/repository/interfaces"
)

// TelemetryUsecase mendefinisikan interface untuk usecase telemetry.
type TelemetryUsecase interface {
	ProcessIncoming(ctx context.Context, payload domain.TelemetryPayload) error
}

type telemetryUsecase struct {
	repo          interfaces.TelemetryRepository
	broadcastChan chan<- []byte
}

// NewTelemetryUsecase membuat instance baru dari TelemetryUsecase.
// Inject repository dan broadcast channel melalui constructor untuk menghindari global variables.
func NewTelemetryUsecase(repo interfaces.TelemetryRepository, broadcastChan chan<- []byte) TelemetryUsecase {
	return &telemetryUsecase{
		repo:          repo,
		broadcastChan: broadcastChan,
	}
}

// ProcessIncoming memproses raw payload dari ESP32, mengevaluasi threshold, 
// menyimpan ke database, dan melakukan broadcast via WebSocket.
func (u *telemetryUsecase) ProcessIncoming(ctx context.Context, payload domain.TelemetryPayload) error {
	// 1. Evaluasi payload untuk mendapatkan TelemetryRecord dengan status threshold
	record := EvaluatePayload(payload)

	// Set waktu pembuatan record (sebelum disimpan ke DB)
	record.CreatedAt = time.Now()

	// 2. Simpan record ke repository (database)
	if err := u.repo.Save(ctx, record); err != nil {
		return fmt.Errorf("failed to save telemetry record: %w", err)
	}

	// 3. Marshal record ke JSON untuk broadcast ke frontend
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry record: %w", err)
	}

	// 4. Kirim ke channel broadcast (non-blocking jika perlu, tapi sesuai blueprint kirim langsung)
	// Pastikan hub sudah menangani channel ini agar tidak deadlock.
	u.broadcastChan <- recordJSON

	// 5. Cek apakah ada status CRITICAL secara keseluruhan
	if record.Status.Overall == domain.StatusCritical {
		alert := domain.Alert{
			DeviceID:    record.DeviceID,
			Status:      record.Status.Overall,
			TriggeredAt: time.Now(),
		}

		// Menentukan AlertType dan Value yang spesifik berdasarkan mana yang kritis.
		if record.Status.Temperature == domain.StatusCritical {
			alert.AlertType = domain.AlertTempCritical
			alert.Value = record.Sensor.Temperature
		} else if record.Status.SpO2 == domain.StatusCritical {
			alert.AlertType = domain.AlertSpO2Critical
			alert.Value = float64(record.Sensor.SpO2)
		}

		// Simpan alert
		if err := u.repo.SaveAlert(ctx, alert); err != nil {
			// Kita return error atau log? Sebaiknya return error dengan context.
			return fmt.Errorf("failed to save critical alert: %w", err)
		}
	}

	return nil
}

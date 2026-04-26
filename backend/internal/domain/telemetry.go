package domain

import "time"

// SensorStatus adalah enum untuk hasil evaluasi threshold.
// Referensi: HEALY_Master_Blueprint.md Section 4.1
type SensorStatus string

const (
	StatusNormal   SensorStatus = "NORMAL"
	StatusWarning  SensorStatus = "WARNING"
	StatusCritical SensorStatus = "CRITICAL"
)

// TelemetryPayload adalah struktur data dari ESP32.
// Harus identik dengan JSON payload yang dikirim oleh iot/src/main.cpp.
type TelemetryPayload struct {
	DeviceID  string     `json:"device_id"`
	Timestamp time.Time  `json:"timestamp"`
	Sensor    SensorData `json:"sensor"`
}

// SensorData berisi pembacaan sensor mentah dari ESP32.
type SensorData struct {
	Temperature float64 `json:"temperature"` // Celsius — MLX90614
	BPM         int     `json:"bpm"`         // Beats per minute — MAX30102
	SpO2        int     `json:"spo2"`        // Percentage — MAX30102
}

// TelemetryRecord adalah data yang disimpan ke DB (sudah dievaluasi).
type TelemetryRecord struct {
	TelemetryPayload
	Status    EvaluatedStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

// EvaluatedStatus berisi hasil evaluasi threshold per sensor.
type EvaluatedStatus struct {
	Temperature SensorStatus `json:"temperature"`
	SpO2        SensorStatus `json:"spo2"`
	Overall     SensorStatus `json:"overall"`
}

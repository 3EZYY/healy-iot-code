package domain

import "time"

// AlertType menentukan jenis alert yang dipicu oleh threshold engine.
type AlertType string

const (
	AlertTempWarning  AlertType = "TEMP_WARNING"
	AlertTempCritical AlertType = "TEMP_CRITICAL"
	AlertSpO2Warning  AlertType = "SPO2_WARNING"
	AlertSpO2Critical AlertType = "SPO2_CRITICAL"
)

// Alert merepresentasikan sebuah entry di tabel alert_logs.
type Alert struct {
	ID          string       `json:"id"`
	DeviceID    string       `json:"device_id"`
	AlertType   AlertType    `json:"alert_type"`
	Value       float64      `json:"value"`
	Status      SensorStatus `json:"status"`
	TriggeredAt time.Time    `json:"triggered_at"`
}

package domain

import "time"

// DeviceSettings berisi konfigurasi threshold per device.
// Default values sesuai Blueprint Section 4.2 dan PRD Section 4.
type DeviceSettings struct {
	DeviceID    string    `json:"device_id"`
	TempWarnMax float64   `json:"temp_warn_max"` // Default: 37.5
	TempCritMax float64   `json:"temp_crit_max"` // Default: 38.5
	SpO2WarnMin int       `json:"spo2_warn_min"` // Default: 94
	SpO2CritMin int       `json:"spo2_crit_min"` // Default: 90
	UpdatedAt   time.Time `json:"updated_at"`
}

// DefaultSettings mengembalikan threshold default sesuai PRD.
func DefaultSettings(deviceID string) DeviceSettings {
	return DeviceSettings{
		DeviceID:    deviceID,
		TempWarnMax: 37.5,
		TempCritMax: 38.5,
		SpO2WarnMin: 94,
		SpO2CritMin: 90,
	}
}

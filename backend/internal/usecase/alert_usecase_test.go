package usecase_test

import (
	"testing"
	"time"

	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/usecase"
)

func TestEvaluateTemperature(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		expected domain.SensorStatus
	}{
		{"Below Normal (Critical)", 36.4, domain.StatusCritical}, // Assuming below min is critical as default case
		{"Normal Lower Bound", 36.5, domain.StatusNormal},
		{"Normal Inside", 37.0, domain.StatusNormal},
		{"Normal Upper Bound", 37.5, domain.StatusNormal},
		{"Warning Lower Bound", 37.6, domain.StatusWarning},
		{"Warning Inside", 38.0, domain.StatusWarning},
		{"Warning Upper Bound", 38.5, domain.StatusWarning},
		{"Critical Above Bound", 38.6, domain.StatusCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usecase.EvaluateTemperature(tt.temp)
			if result != tt.expected {
				t.Errorf("EvaluateTemperature(%.1f) = %v; want %v", tt.temp, result, tt.expected)
			}
		})
	}
}

func TestEvaluateSpO2(t *testing.T) {
	tests := []struct {
		name     string
		spo2     int
		expected domain.SensorStatus
	}{
		{"Normal Upper Bound", 100, domain.StatusNormal},
		{"Normal Lower Bound", 95, domain.StatusNormal},
		{"Warning Upper Bound", 94, domain.StatusWarning},
		{"Warning Lower Bound", 91, domain.StatusWarning},
		{"Critical Below Bound", 90, domain.StatusCritical},
		{"Critical Very Low", 85, domain.StatusCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usecase.EvaluateSpO2(tt.spo2)
			if result != tt.expected {
				t.Errorf("EvaluateSpO2(%d) = %v; want %v", tt.spo2, result, tt.expected)
			}
		})
	}
}

func TestEvaluateOverall(t *testing.T) {
	tests := []struct {
		name       string
		tempStatus domain.SensorStatus
		spo2Status domain.SensorStatus
		expected   domain.SensorStatus
	}{
		{"Both Normal", domain.StatusNormal, domain.StatusNormal, domain.StatusNormal},
		{"Temp Warning", domain.StatusWarning, domain.StatusNormal, domain.StatusWarning},
		{"SpO2 Warning", domain.StatusNormal, domain.StatusWarning, domain.StatusWarning},
		{"Both Warning", domain.StatusWarning, domain.StatusWarning, domain.StatusWarning},
		{"Temp Critical", domain.StatusCritical, domain.StatusNormal, domain.StatusCritical},
		{"SpO2 Critical", domain.StatusNormal, domain.StatusCritical, domain.StatusCritical},
		{"Critical and Warning", domain.StatusCritical, domain.StatusWarning, domain.StatusCritical},
		{"Both Critical", domain.StatusCritical, domain.StatusCritical, domain.StatusCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usecase.EvaluateOverall(tt.tempStatus, tt.spo2Status)
			if result != tt.expected {
				t.Errorf("EvaluateOverall(%v, %v) = %v; want %v", tt.tempStatus, tt.spo2Status, result, tt.expected)
			}
		})
	}
}

func TestEvaluatePayload(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		payload  domain.TelemetryPayload
		expected domain.EvaluatedStatus
	}{
		{
			name: "All Normal",
			payload: domain.TelemetryPayload{
				DeviceID:  "healy-001",
				Timestamp: now,
				Sensor: domain.SensorData{
					Temperature: 37.0,
					BPM:         75,
					SpO2:        98,
				},
			},
			expected: domain.EvaluatedStatus{
				Temperature: domain.StatusNormal,
				SpO2:        domain.StatusNormal,
				Overall:     domain.StatusNormal,
			},
		},
		{
			name: "Mixed Warning and Critical",
			payload: domain.TelemetryPayload{
				DeviceID:  "healy-001",
				Timestamp: now,
				Sensor: domain.SensorData{
					Temperature: 38.0, // Warning
					BPM:         120,
					SpO2:        89,   // Critical
				},
			},
			expected: domain.EvaluatedStatus{
				Temperature: domain.StatusWarning,
				SpO2:        domain.StatusCritical,
				Overall:     domain.StatusCritical,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usecase.EvaluatePayload(tt.payload)

			if result.Status != tt.expected {
				t.Errorf("EvaluatePayload() Status = %v; want %v", result.Status, tt.expected)
			}
			if result.DeviceID != tt.payload.DeviceID {
				t.Errorf("EvaluatePayload() DeviceID = %s; want %s", result.DeviceID, tt.payload.DeviceID)
			}
			if result.Timestamp != tt.payload.Timestamp {
				t.Errorf("EvaluatePayload() Timestamp = %v; want %v", result.Timestamp, tt.payload.Timestamp)
			}
			if result.Sensor != tt.payload.Sensor {
				t.Errorf("EvaluatePayload() Sensor = %v; want %v", result.Sensor, tt.payload.Sensor)
			}
		})
	}
}

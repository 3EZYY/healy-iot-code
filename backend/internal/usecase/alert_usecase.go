package usecase

import "github.com/rafif/healy-backend/internal/domain"

// Threshold defaults — overrideable via device_settings table
const (
	TempNormalMin float64 = 36.5
	TempNormalMax float64 = 37.5
	TempWarnMax   float64 = 38.5 // Above this = CRITICAL

	SpO2NormalMin int = 95
	SpO2WarnMin   int = 91 // Below this = CRITICAL
)

// EvaluateTemperature mengembalikan status berdasarkan rentang suhu.
func EvaluateTemperature(temp float64) domain.SensorStatus {
	switch {
	case temp >= TempNormalMin && temp <= TempNormalMax:
		return domain.StatusNormal
	case temp > TempNormalMax && temp <= TempWarnMax:
		return domain.StatusWarning
	default:
		return domain.StatusCritical
	}
}

// EvaluateSpO2 mengembalikan status berdasarkan kadar oksigen dalam darah.
func EvaluateSpO2(spo2 int) domain.SensorStatus {
	switch {
	case spo2 >= SpO2NormalMin:
		return domain.StatusNormal
	case spo2 >= SpO2WarnMin:
		return domain.StatusWarning
	default:
		return domain.StatusCritical
	}
}

// EvaluateOverall menentukan status keseluruhan dari status individual sensor.
func EvaluateOverall(tempStatus, spo2Status domain.SensorStatus) domain.SensorStatus {
	if tempStatus == domain.StatusCritical || spo2Status == domain.StatusCritical {
		return domain.StatusCritical
	}
	if tempStatus == domain.StatusWarning || spo2Status == domain.StatusWarning {
		return domain.StatusWarning
	}
	return domain.StatusNormal
}

// EvaluatePayload adalah entry point utama — menerima raw payload dari ESP32
// dan mengembalikan TelemetryRecord yang siap disimpan ke DB
func EvaluatePayload(payload domain.TelemetryPayload) domain.TelemetryRecord {
	tempStatus := EvaluateTemperature(payload.Sensor.Temperature)
	spo2Status := EvaluateSpO2(payload.Sensor.SpO2)

	return domain.TelemetryRecord{
		TelemetryPayload: payload,
		Status: domain.EvaluatedStatus{
			Temperature: tempStatus,
			SpO2:        spo2Status,
			Overall:     EvaluateOverall(tempStatus, spo2Status),
		},
	}
}

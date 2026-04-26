package interfaces

import (
	"context"

	"github.com/rafif/healy-backend/internal/domain"
)

type TelemetryRepository interface {
	Save(ctx context.Context, record domain.TelemetryRecord) error
	SaveAlert(ctx context.Context, alert domain.Alert) error
	GetHistory(ctx context.Context, deviceID string, hours int) ([]domain.TelemetryRecord, error)
	GetLatest(ctx context.Context, deviceID string) (domain.TelemetryRecord, error)
}

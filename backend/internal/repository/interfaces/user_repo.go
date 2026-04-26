package interfaces

import (
	"context"

	"github.com/rafif/healy-backend/internal/domain"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

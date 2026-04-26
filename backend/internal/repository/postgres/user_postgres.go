package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/repository/interfaces"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	const query = `
		SELECT id, username, password, created_at
		FROM public.users
		WHERE username = $1
	`
	
	var user domain.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

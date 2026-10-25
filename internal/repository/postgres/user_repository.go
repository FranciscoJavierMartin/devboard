package postgres

import (
	"context"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	generated "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRespository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewUserRespository(pool *pgxpool.Pool) *UserRespository {
	return &UserRespository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *UserRespository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	})

	if err != nil {
		return fmt.Errorf("inserting user on DB: %w", err)
	}

	user.ID = result.ID.String()
	user.CreatedAt = result.CreatedAt.Time
	user.UpdatedAt = result.UpdatedAt.Time

	return nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	generated "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres/generated"
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

func (r *UserRespository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, domain.ErrNotFound)
	}

	user, err := r.queries.GetUserByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving user %s: %w", id, domain.ErrInternal)
	}

	return toDomainUser(user), nil
}

func (r *UserRespository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", email, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving user %s: %w", email, domain.ErrInternal)
	}

	return toDomainUser(user), nil
}

func toDomainUser(u generated.User) *domain.User {
	return &domain.User{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

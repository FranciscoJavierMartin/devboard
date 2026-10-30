package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	generated "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewWorkspaceRepository(pool *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *WorkspaceRepository) CreateWorkspace(ctx context.Context, ws *domain.Workspace) error {
	var createdBy pgtype.UUID

	if err := createdBy.Scan(ws.CreatedBy); err != nil {
		return fmt.Errorf("invalid created_by: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateWorkspace(ctx, generated.CreateWorkspaceParams{
		Name:        ws.Name,
		Slug:        ws.Slug,
		Description: ws.Description,
		CreatedBy:   createdBy,
	})

	if err != nil {
		return fmt.Errorf("inserting workspace on DB: %w", err)
	}

	ws.ID = result.ID.String()
	ws.CreatedAt = result.CreatedAt.Time
	ws.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r *WorkspaceRepository) GetWorkspaceByID(ctx context.Context, id string) (*domain.Workspace, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, domain.ErrNotFound)
	}

	workspace, err := r.queries.GetWorkspaceByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("workspace %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving workspace %s: %w", id, domain.ErrNotFound)
	}

	return toDomainWorkspace(workspace), nil
}

func (r *WorkspaceRepository) GetWorkspaceBySlug(ctx context.Context, slug string) (*domain.Workspace, error) {
	workspace, err := r.queries.GetWorkspaceBySlug(ctx, slug)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("workspace %s: %w", slug, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving workspace %s: %w", slug, err)
	}

	return toDomainWorkspace(workspace), nil
}

func (r *WorkspaceRepository) ListWorkspacesByUser(ctx context.Context, userID string) ([]*domain.Workspace, error) {
	var pgUserID pgtype.UUID

	if err := pgUserID.Scan(userID); err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", domain.ErrInvalidInput)
	}
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, member *domain.WorkspaceMember) error {
	var wsID, userID pgtype.UUID

	if err := wsID.Scan(member.WorkspaceID); err != nil {
		return fmt.Errorf("invalid workspace_id: %w", domain.ErrInvalidInput)
	}

	if err := userID.Scan(member.UserID); err != nil {
		return fmt.Errorf("invalid user_id: %w", domain.ErrInvalidInput)
	}
}

func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID string) (*domain.WorkspaceMember, error) {
	var pgWsID, pgUserID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("invalid workspace_id: %w", domain.ErrNotFound)
	}

	if err := pgUserID.Scan(userID); err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", domain.ErrNotFound)
	}
}

func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID string) ([]*domain.WorkspaceMember, error) {
	var pgWsID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("invalid workspace_id: %w", domain.ErrInvalidInput)
	}
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	var pgWsID, pgUserID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return fmt.Errorf("invalid workspace_id: %w", domain.ErrInvalidInput)
	}

	if err := pgUserID.Scan(userID); err != nil {
		return fmt.Errorf("invalid user_id: %w", domain.ErrInvalidInput)
	}
}

func toDomainWorkspace(w generated.Workspace) *domain.Workspace {
	return &domain.Workspace{
		ID:          w.ID.String(),
		Name:        w.Name,
		Slug:        w.Slug,
		Description: w.Description,
		CreatedBy:   w.CreatedBy.String(),
		CreatedAt:   w.CreatedAt.Time,
		UpdatedAt:   w.UpdatedAt.Time,
	}
}

func toDomainMember(m generated.WorkspaceMember) *domain.WorkspaceMember {
	return &domain.WorkspaceMember{
		ID:          m.ID.String(),
		WorkspaceID: m.WorkspaceID.String(),
		UserID:      m.UserID.String(),
		Role:        domain.WorkspaceRole(m.Role),
		JoinedAt:    m.JoinedAt.Time,
	}
}

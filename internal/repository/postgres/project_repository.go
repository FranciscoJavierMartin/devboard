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

type ProjectRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	var wsID, createdBy pgtype.UUID

	if err := wsID.Scan(project.WorkspaceID); err != nil {
		return fmt.Errorf("invalid workspace_id: %w", domain.ErrInvalidInput)
	}

	if err := createdBy.Scan(project.CreatedBy); err != nil {
		return fmt.Errorf("invalid created_by: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateProject(ctx, generated.CreateProjectParams{
		WorkspaceID: wsID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   createdBy,
	})

	if err != nil {
		return fmt.Errorf("inserting project into DB: %w", err)
	}

	project.ID = result.ID.String()
	project.CreatedAt = result.CreatedAt.Time
	project.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, domain.ErrNotFound)
	}

	project, err := r.queries.GetProjectByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("project %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving project %s: %w", id, err)
	}

	return toDomainProject(project), nil
}

func (r *ProjectRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*domain.Project, error) {
	var pgWsID pgtype.UUID

	if err := pgWsID.Scan(workspaceID); err != nil {
		return nil, fmt.Errorf("invalid workspace_id: %w", domain.ErrInvalidInput)
	}
}

func toDomainProject(p generated.Project) *domain.Project {
	return &domain.Project{
		ID:          p.ID.String(),
		WorkspaceID: p.WorkspaceID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedBy:   p.CreatedBy.String(),
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}

package postgres

import (
	"context"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	generated "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRespository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRespository {
	return &TaskRespository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *TaskRespository) Create(ctx context.Context, task *domain.Task) error {
	var projetID pgtype.UUID
	if err := projetID.Scan(task.ProjectID); err != nil {
		return fmt.Errorf("invalid project_id: %w", domain.ErrInvalidInput)
	}

	var createdBy pgtype.UUID
	if err := createdBy.Scan(task.CreatedBy); err != nil {
		return fmt.Errorf("invalid created_by: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateTask(ctx, generated.CreateTaskParams{
		ProjectID:   projetID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedBy:   createdBy,
	})

	if err != nil {
		return fmt.Errorf("inserting task into DB: %w", err)
	}

	task.ID = result.ID.String()
	task.CreatedAt = result.CreatedAt.Time
	task.UpdatedAt = result.UpdatedAt.Time

	return nil
}

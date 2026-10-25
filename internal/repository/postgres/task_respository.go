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

func (r *TaskRespository) Update(ctx context.Context, task *domain.Task) error {
	var id pgtype.UUID
	if err := id.Scan(task.ID); err != nil {
		return fmt.Errorf("invalid id: %w", domain.ErrInvalidInput)
	}

	var assigneeID pgtype.UUID

	if task.AssigneeID != nil {
		if err := assigneeID.Scan(task.AssigneeID); err != nil {
			return fmt.Errorf("invalid assigned_id: %w", domain.ErrInvalidInput)
		}
	}

	_, err := r.queries.UpdateTask(ctx, generated.UpdateTaskParams{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		AssigneeID:  assigneeID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("task %s: %w", task.ID, domain.ErrNotFound)
		}
		return fmt.Errorf("updating task %s: %w", task.ID, err)
	}

	return nil
}

func (r *TaskRespository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var taskID pgtype.UUID

	if err := taskID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", domain.ErrInvalidInput)
	}

	task, err := r.queries.GetTaskByID(ctx, taskID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving task: %s: %w", id, err)
	}

	return toDomainTask(task), nil
}

func toDomainTask(t generated.Task) *domain.Task {
	task := &domain.Task{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      domain.TaskStatus(t.Status),
		CreatedBy:   t.CreatedBy.String(),
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}

	if t.AssigneeID.Valid {
		id := t.AssigneeID.String()
		task.AssigneeID = &id
	}

	return task
}

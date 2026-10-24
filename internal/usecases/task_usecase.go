package usecases

import (
	"context"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

type TaskCommandRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
}

type TaskQueryRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	ListByProject(ctx context.Context, projectId string) ([]*domain.Task, error)
}

type TaskUserCase struct {
	commands TaskCommandRepository
	queries  TaskQueryRepository
}

func NewTaskUseCase(commands TaskCommandRepository, queries TaskQueryRepository) *TaskUserCase {
	return &TaskUserCase{
		commands: commands,
		queries:  queries,
	}
}

func (uc *TaskUserCase) CreateTask(ctx context.Context, projectID, title, createdBy string) (*domain.Task, error) {
	task := domain.NewTask(projectID, title, createdBy)

	if err := uc.commands.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}

	return task, nil
}

func (uc *TaskUserCase) ListProjectTasks(ctx context.Context, projectID string) ([]*domain.Task, error) {
	return uc.queries.ListByProject(ctx, projectID)
}

func (uc *TaskUserCase) AssignTask(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	task, err := uc.queries.GetByID(ctx, taskID)

	if err != nil {
		return nil, fmt.Errorf("getting task %s: %w", taskID, err)
	}

	if err := task.AssignTo(userID); err != nil {
		return nil, err
	}

	if err := uc.commands.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}

	return task, nil
}

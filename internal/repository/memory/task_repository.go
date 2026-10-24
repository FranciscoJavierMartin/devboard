package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *TaskRepository) Create(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = fmt.Sprintf("task_%d", len(r.tasks)+1)
	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) Update(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) GetByID(_ context.Context, taskID string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RLocker().Unlock()

	task, ok := r.tasks[taskID]

	if !ok {
		return nil, fmt.Errorf("task %s: %w", taskID, domain.ErrNotFound)
	}

	return task, nil
}

func (r *TaskRepository) ListByProject(_ context.Context, projectID string) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RLocker().Unlock()

	var tasks []*domain.Task

	for _, task := range r.tasks {

		if task.ProjectID == projectID {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

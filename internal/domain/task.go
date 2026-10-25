package domain

import (
	"errors"
	"time"
)

// Task represent a task inside the board
type Task struct {
	ID          string
	ProjectID   string
	Title       string
	Description string
	Status      TaskStatus
	AssigneeID  *string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(projectID, title, description, createdBy string) *Task {
	now := time.Now()

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		CreatedBy:   createdBy,
		Status:      TaskStatusTodo,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

var (
	ErrTaskAlreadyDone   = errors.New("Task is already done")
	ErrTaskNotAssignable = errors.New("Task cannot be assigned to that state")
)

// AssignTo assign task to user
func (task *Task) AssignTo(userID string) error {
	if task.Status == TaskStatusDone {
		return ErrTaskNotAssignable
	}

	task.AssigneeID = &userID
	task.UpdatedAt = time.Now()
	return nil
}

// MarkAsDone mark task as done
func (task *Task) MarkAsDone() error {
	if task.Status == TaskStatusDone {
		return ErrTaskAlreadyDone
	}

	task.Status = TaskStatusDone
	task.UpdatedAt = time.Now()
	return nil
}

// MoveToStatus change task status
func (task *Task) MoveToStatus(status TaskStatus) error {
	if !status.IsValid() {
		return errors.New("New status is invalid")
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	return nil
}

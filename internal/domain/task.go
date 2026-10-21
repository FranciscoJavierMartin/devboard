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

func NewTask(projectID, title, createdBy string) *Task {
	now := time.Now()

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		CreatedBy:   createdBy,
		Status:      TaskStatusTodo,
		Description: "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

var (
	ErrTaskAlreadyDone   = errors.New("Task is already done")
	ErrTaskNotAssignable = errors.New("Task cannot be assigned to that state")
)

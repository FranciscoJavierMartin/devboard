// Package domain task status
package domain

import (
	"encoding/json"
	"fmt"
)

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	// TaskStatusTodo No initialized
	TaskStatusTodo TaskStatus = "todo"
	// TaskStatusInProgress In progress
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusDone Done
	TaskStatusDone TaskStatus = "done"
)

// IsValid verify that status is one of the allowed values
func (status TaskStatus) IsValid() bool {
	switch status {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// UnmarshalJSON validate status at desialize
func (status *TaskStatus) UnmarshalJSON(data []byte) error {
	var raw string

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	statusTask := TaskStatus(raw)

	if !statusTask.IsValid() {
		return fmt.Errorf("invalid status: %q", raw)
	}

	*status = statusTask
	return nil
}

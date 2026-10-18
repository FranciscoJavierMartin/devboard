package domain

import (
	"encoding/json"
	"testing"
)

type CreateTaskRequest struct {
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

func TestCreateTaskRequestWithValidStatus(test *testing.T) {
	body := []byte(`{"title": "My task", "status": "todo"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)

	if err != nil {
		test.Fatalf("Unexpected error, received: %v", err)
	}

	if req.Status != TaskStatusTodo {
		test.Fatalf("Unexected status, received: %s", req.Status)
	}
}

func TestCreateTaskRequestWithInvalidStatus(test *testing.T) {
	body := []byte(`{"title": "My task", "status": "invalid"}`)

	var req CreateTaskRequest

	err := json.Unmarshal(body, &req)

	if err == nil {
		test.Fatal("Expected error, but received nil")
	}

	expected := `invalid status: "invalid"`

	if err.Error() != expected {
		test.Fatalf("Unexected error %q, received: %q", expected, err.Error())
	}
}

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	"github.com/FranciscoJavierMartin/devboard/internal/usecases"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
)

type TaskHandler struct {
	usercase  *usecases.TaskUserCase
	validator *validator.Validator
	logger    *slog.Logger
}

func NewTaskHandler(uc *usecases.TaskUserCase, v *validator.Validator, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		usercase:  uc,
		validator: v,
		logger:    logger,
	}
}

type createTaskRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	Title     string `json:"title" validate:"required,min=3,max=200"`
}

type assignTaskRequest struct {
	AssigneeID string `json:"assignee_id" validate:"required,uuid4"`
}

type taskResponse struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Status      domain.TaskStatus `json:"status"`
	AssigneeID  *string           `json:"assignee_id,omitempty"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// toTaskResponse transform entity to response DTO
func toTaskResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeID:  t.AssigneeID,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// Create method to create tasks
//
// @Summary	Create task
// @Description Create a new task
// @Tags  tasks
// @Accept json
// @Produce json
// @Param request body	createTaskRequest	true "Task data"
// @Success 201 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error writing validation response",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	task, err := h.usercase.CreateTask(r.Context(), req.ProjectID, req.Title, "system")
	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing JSON response",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// Get method to retrieve task by ID
//
// @Summary Retrieve task by ID
// @Description Retrieve a task by its ID
// @Tags  tasks
// @Produce json
// @Param id path	string	true "task ID"
// @Success 200 {object} taskResponse
// @Failure 404 {object} ProblemDetails
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.usercase.GetTask(r.Context(), id)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing JSON response",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// ListByProject method to retrieve tasks belonging to a project
//
// @Summary	List all tasks from a project
// @Description Returns all tasks that belongs to a project
// @Tags  tasks
// @Produce json
// @Param id path	string	true "Project ID"
// @Success 200 {object} PaginatedResponse[taskResponse]
// @Failure 404 {object} ProblemDetails
// @Router /projects/{id}/tasks [get]
func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	tasks, err := h.usercase.ListProjectTasks(r.Context(), projectID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	response := make([]taskResponse, len(tasks))

	for i, task := range tasks {
		response[i] = toTaskResponse(task)
	}

	if err := RespondJSON(w, http.StatusOK, PaginatedResponse[taskResponse]{
		Data:    response,
		HasMore: false,
	}); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error getting JSON response to get tasks",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

// Assign method to assign a task to a user
//
// @Summary	Assign task
// @Description Assign task to a team member
// @Tags  tasks
// @Accept json
// @Produce json
// @Param id 	path	string		true 	"Task ID"
// @Param request body	assignTaskRequest	true "User ID"
// @Success 200 {object} taskResponse
// @Failure 400 {object} ProblemDetails
// @Router /tasks/{id}/assign [put]
func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req assignTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, r, h.logger, domain.ErrInvalidInput)
	}

	if errs := h.validator.Validate(req); errs != nil {
		if err := RespondValidationError(w, errs); err != nil {
			h.logger.ErrorContext(
				r.Context(),
				"error writing validation response",
				"error", err,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
		return
	}

	task, err := h.usercase.AssignTask(r.Context(), id, req.AssigneeID)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusOK, toTaskResponse(task)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing JSON response",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}

}

// UpdateStatus

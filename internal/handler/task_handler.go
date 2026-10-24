package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

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

type createTaskResponse struct {
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
}

// Create method to create tasks
//
// @Summary	Create task
// @Description Create a new task
// @Tags  tasks
// @Accept json
// @Produce json
// @Param request body	createTaskRequest	true "Task data"
// @Success 201 {object} createTaskResponse
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

	if err := respondJSON(w, http.StatusCreated, createTaskResponse{
		ProjectID: task.ProjectID,
		Title:     task.Title,
	}); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing JSON response",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

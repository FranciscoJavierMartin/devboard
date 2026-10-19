package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
)

type errorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

// ProblemDetails struct standarized by rfc 7807
type ProblemDetails struct {
	Type     string                      `json:"type"`
	Title    string                      `json:"title"`
	Status   int                         `json:"status"`
	Detail   string                      `json:"detail"`
	Instance string                      `json:"instance"`
	Errors   []validator.ValidationError `json:"errors,omitempty"`
}

const problemBaseURL = "http://localhost:8080"

// RespondError centralize error response
func RespondError(writer http.ResponseWriter, request *http.Request, logger *slog.Logger, err error) {
	var problem ProblemDetails
	problem.Instance = request.URL.Path

	switch {
	case errors.Is(err, domain.ErrNotFound):
		problem.Type = problemBaseURL + "/not-found"
		problem.Title = "Resource not found"
		problem.Status = http.StatusNotFound
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		problem.Type = problemBaseURL + "/conflict"
		problem.Title = "Resource already exists"
		problem.Status = http.StatusConflict
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		problem.Type = problemBaseURL + "/unauthorized"
		problem.Title = "Unauthorized"
		problem.Status = http.StatusUnauthorized
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		problem.Type = problemBaseURL + "/forbidden"
		problem.Title = "Forbidden"
		problem.Status = http.StatusForbidden
		problem.Detail = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		problem.Type = problemBaseURL + "/bad-request"
		problem.Title = "Bad request"
		problem.Status = http.StatusBadRequest
		problem.Detail = err.Error()
	default:
		problem.Type = problemBaseURL + "/internal"
		problem.Title = "Internal server error"
		problem.Status = http.StatusInternalServerError
		problem.Detail = "Something unexpected happened"
		logger.ErrorContext(
			request.Context(),
			"unclassified error",
			slog.Any("error", err),
			slog.String("path", request.URL.Path),
		)
	}

	writer.Header().Set("Content-Type", "application/problem+json")
	writer.WriteHeader(problem.Status)

	if encodeErr := json.NewEncoder(writer).Encode(problem); encodeErr != nil {
		logger.Error(
			"error writing problem details",
			slog.String("error", encodeErr.Error()),
		)
	}
}

// RespondValidationError responds validation errors
func RespondValidationError(writer http.ResponseWriter, errs []validator.ValidationError) {
	respondJSON(writer, http.StatusBadRequest, errorResponse{
		Error:   "invalid input data",
		Details: errs,
	})
}

// ResponseJSON send a successful response in JSON format
func ResponseJSON(writer http.ResponseWriter, status int, data any) {
	respondJSON(writer, status, data)
}

func respondJSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(data)

	if err != nil {
		logger := logger.New(logger.DefaultConfig())
		logger.Debug("Error creating Scalar page: ", err.Error(), "")
	}
}

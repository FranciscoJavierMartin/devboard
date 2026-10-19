package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
)

type errorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

// Centralized RespondError
func RespondError(writer http.ResponseWriter, request *http.Request, logger *slog.Logger, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	default:
		status = http.StatusInternalServerError
		message = "internal server error"
		logger.ErrorContext(
			request.Context(),
			"unclassified error",
			slog.Any("error", err),
			slog.String("path", request.URL.Path),
		)
	}

	respondJSON(writer, status, errorResponse{Error: message})
}

func RespondValidationError(writer http.ResponseWriter, errs []validator.ValidationError) {
	respondJSON(writer, http.StatusBadRequest, errorResponse{
		Error:   "invalid input data",
		Details: errs,
	})
}

func ResponseJSON(writer http.ResponseWriter, status int, data any) {
	respondJSON(writer, status, data)
}

func respondJSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(data)
}

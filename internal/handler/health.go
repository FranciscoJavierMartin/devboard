// Package handler health
package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler struct
type HealthHandler struct{}

// NewHealthHandler contructor
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// healthResponse struct
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// HealthCheck verify that server is running
//
// @Summary	Health check
// @Description Verify that server is running and accepting requests
// @Tags  system
// @Produce json
// @Success 200 {object} healthResponse
// @Router /health [get]
func (handle *HealthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	resp := healthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return
	}
}

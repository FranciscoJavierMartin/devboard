package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (handle *HealthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

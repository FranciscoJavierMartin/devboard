package handler

import (
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Cursor string
	Limit  int
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

// ParsePagination read cursor ant limit from query string
func ParsePagination(request *http.Request) PaginationParams {
	cursor := request.URL.Query().Get("cursor")
	limit := defaultLimit

	if limitQuery := request.URL.Query().Get("limit"); limitQuery != "" {
		if parsed, err := strconv.Atoi(limitQuery); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return PaginationParams{
		Cursor: cursor,
		Limit:  limit,
	}
}

// PaginatedResponse is the standard format for any paginated list
type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

package middlewares

import (
	"net/http"
)

// Middleware is the signature for any middleware
type Middleware func(http.Handler) http.Handler

// Chain applies multiples middlewares to a handler in reverse order
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

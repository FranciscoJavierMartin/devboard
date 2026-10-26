package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FranciscoJavierMartin/devboard/internal/auth"
	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	"github.com/FranciscoJavierMartin/devboard/internal/handler"
)

type contextKey string

const (
	claimsKey contextKey = "auth_claims"
)

func RequireAuth(jwtManager *auth.JWTManager, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractBearerToken(r)

			if err != nil {
				handler.RespondError(w, r, logger, domain.ErrUnauthorized)
				return
			}

			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				handler.RespondError(w, r, logger, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", domain.ErrUnauthorized
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", domain.ErrUnauthorized
	}

	token := strings.TrimSpace(parts[1])

	if token == "" {
		return "", domain.ErrUnauthorized
	}

	return token, nil
}

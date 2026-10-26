package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	"github.com/FranciscoJavierMartin/devboard/internal/usecases"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
)

type AuthHandler struct {
	usecase   *usecases.AuthUseCase
	validator *validator.Validator
	logger    *slog.Logger
}

func NewAuthHandler(uc *usecases.AuthUseCase, v *validator.Validator, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		usecase:   uc,
		validator: v,
		logger:    logger,
	}
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Register endpoint
//
// @Summary     Register user
// @Description Create new user and returns access token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     registerRequest true "register data"
// @Success     201     {object} tokenResponse
// @Failure     400     {object} ProblemDetails
// @Failure     409     {object} ProblemDetails
// @Router      /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

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

	tokens, err := h.usecase.Register(r.Context(), req.Email, req.Name, req.Password)

	if err != nil {
		RespondError(w, r, h.logger, err)
		return
	}

	if err := RespondJSON(w, http.StatusCreated, toTokenResponse(tokens)); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"error writing JSON response",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	}
}

func toTokenResponse(tokens *usecases.AuthToken) tokenResponse {
	return tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(tokens.ExpiresIn.Seconds()),
	}
}

package handler

import (
	"log/slog"

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

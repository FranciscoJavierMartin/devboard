package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FranciscoJavierMartin/devboard/internal/auth"
	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

type authUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

type AuthUseCase struct {
	userRepo             authUserRepository
	jwtManager           *auth.JWTManager
	tokenStore           *auth.TokenStore
	refreshTokenDuration time.Duration
}

func NewAuthUseCase(
	userRepo authUserRepository,
	jwtManager *auth.JWTManager,
	tokenStore *auth.TokenStore,
	refreshTokenDuration time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:             userRepo,
		jwtManager:           jwtManager,
		tokenStore:           tokenStore,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, name, password string) (*AuthToken, error) {
	passwordHash, err := auth.HashPassword(password)

	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := domain.NewUser(email, name, passwordHash)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *AuthUseCase) generateTokenPair(ctx context.Context, userID, email, role string) (*AuthToken, error) {
	accessToken, err := uc.jwtManager.GenerateAccessToken(userID, email, role)

	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(userID, email, role)

	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	refreshClaims, err := uc.jwtManager.VerifyToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("extracting claims from refresh token: %w", err)
	}

	if err := uc.tokenStore.StoreRefreshToken(
		ctx,
		refreshClaims.ID,
		userID,
		uc.refreshTokenDuration,
	); err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * time.Minute,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*AuthToken, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("retrieving user: %w", err)
	}

	if err := auth.VerifyPassword(password, user.PasswordHash); err != nil {
		if errors.Is(err, auth.ErrPasswordMissmatch) {
			return nil, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("verifying password: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, user.ID, user.Email, "member")

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthToken, error) {
	claims, err := uc.jwtManager.VerifyToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", domain.ErrUnauthorized)
	}

	_, err = uc.tokenStore.ValidateRefreshToken(ctx, claims.ID)

	if err != nil {

		if errors.Is(err, auth.ErrTokenRevoked) || errors.Is(err, auth.ErrTokenNotFound) {
			return nil, fmt.Errorf("invalid refresh token: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("validating refresh token: %w", err)
	}

	if err := uc.tokenStore.RevokeRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return nil, fmt.Errorf("revoking previous token: %w", err)
	}

	tokens, err := uc.generateTokenPair(ctx, claims.ID, claims.Email, claims.Role)

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := uc.jwtManager.VerifyToken(refreshToken)

	if err != nil {
		return nil
	}

	if err := uc.tokenStore.RevokeRefreshToken(ctx, claims.ID, uc.refreshTokenDuration); err != nil {
		return fmt.Errorf("revoking token on logout: %w", err)
	}

	return nil
}

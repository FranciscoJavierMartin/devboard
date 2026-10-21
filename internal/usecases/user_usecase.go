package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

// UserRepository define that Use Case needs
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// UserUseCase organize business about user
type UserUseCase struct {
	repo UserRepository
}

// NewUserUseCase constructor
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

// CreateUser create user and validate it
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name, passwordHash string) (*domain.User, error) {
	existing, err := uc.repo.GetByEmail(ctx, email)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("verifying existing email: %w", err)
	}

	if existing != nil {
		return nil, fmt.Errorf("email %s: %w", email, domain.ErrAlreadyExists)
	}

	user := domain.NewUser(email, name, passwordHash)

	if !user.IsValid() {
		return nil, domain.ErrInvalidInput
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}

	return user, nil
}

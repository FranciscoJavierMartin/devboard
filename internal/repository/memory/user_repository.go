package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

// UserRepository is an adapter which implements UserRepository
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("user_%d", len(r.users)+1)
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) GetByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]

	if !ok {
		return nil, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
	}

	return user, nil
}

// GetByEmail get user by email
func (r *UserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil

		}
	}

	return nil, fmt.Errorf("email: %s: %w", email, domain.ErrNotFound)
}

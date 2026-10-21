package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
)

func TestUserRepository_Create(test *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "John",
		Email: "john@test.com",
	}

	err := repo.Create(ctx, user)

	if err != nil {
		test.Fatalf("Create() returns an unexpected error: %v", err)
	}

	if user.ID == "" {
		test.Fatalf("Create() unassigned ID to user")
	}

	if user.ID != "user_1" {
		test.Errorf("Create() assign ID: %q; expected %q", user.ID, "user_1")
	}
}

func TestUserRepository_GetByID_NotFound(test *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "user_999")

	if err == nil {
		test.Fatalf("GetByID() should return an error")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		test.Errorf("GetByID() returns %v: expected ErrNotFound", err)
	}
}

// Test _GetByID
func TestUserRepository_GetByID(test *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "John",
		Email: "john@test.com",
	}

	if err := repo.Create(ctx, user); err != nil {
		test.Fatalf("User cannot be created: %v", err)
	}

	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		test.Fatalf("GetByID() returns an unexpected error: %v", err)
	}

	if found.ID != user.ID {
		test.Errorf(
			"GetByID() returns ID %q; was expected %q",
			found.ID,
			user.ID,
		)
	}

	if found.Email != user.Email {
		test.Errorf(
			"GetByID() returns email %q; was expected %q",
			found.Email,
			user.Email,
		)
	}
}

// Test _GetByEmail
func TestUserRepository_GetByEmail(test *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:  "John",
		Email: "john@test.com",
	}

	if err := repo.Create(ctx, user); err != nil {
		test.Fatalf("User cannot be created: %v", err)
	}

	found, err := repo.GetByEmail(ctx, "john@test.com")
	if err != nil {
		test.Fatalf(
			"GetByEmail() returns an unexpected error: %v",
			err,
		)
	}

	if found.ID != user.ID {
		test.Errorf(
			"GetByEmail() returns ID %q; was expected %q",
			found.ID,
			user.ID,
		)
	}
}

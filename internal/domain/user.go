package domain

import "time"

// User represent the user
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, name, passwordHash string) *User {
	now := time.Now()

	return &User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsValid verify that email and name are not empty
func (user *User) IsValid() bool {
	return user.Email != "" && user.Name != ""
}

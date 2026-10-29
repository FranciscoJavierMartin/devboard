package domain

import (
	"strings"
	"time"
	"unicode"
)

type WorkspaceRole string

const (
	WorkspaceRoleAdmin  WorkspaceRole = "admin"
	WorkspaceRoleMember WorkspaceRole = "member"
	WorkspaceRoleViewer WorkspaceRole = "viewer"
)

func (r WorkspaceRole) IsValid() bool {
	var result bool

	switch r {
	case WorkspaceRoleAdmin, WorkspaceRoleMember, WorkspaceRoleViewer:
		result = true
	}

	return result
}

type Workspace struct {
	ID          string
	Name        string
	Slug        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkspaceMember struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        WorkspaceRole
	JoinedAt    time.Time
}

func NewWorkspace(name, description, createdBy string) *Workspace {
	now := time.Now()

	return &Workspace{
		Name:        name,
		Slug:        GenerateSlug(name),
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	var result strings.Builder

	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	cleaned := strings.Join(strings.FieldsFunc(result.String(), func(r rune) bool {
		return r == '-'
	}), "-")

	return cleaned
}

func NewWorkspaceMember(workspaceID, userID string, role WorkspaceRole) *WorkspaceMember {
	return &WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now(),
	}
}

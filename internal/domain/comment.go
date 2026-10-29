package domain

import (
	"errors"
	"time"
)

var ErrCommentEmpty = errors.New("comment cannot be empty")

type Comment struct {
	ID        string
	TaskID    string
	AuthorID  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComment(taskID, authorID, content string) (*Comment, error) {
	now := time.Now()

	if content == "" {
		return nil, ErrCommentEmpty
	}

	return &Comment{
		TaskID:    taskID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

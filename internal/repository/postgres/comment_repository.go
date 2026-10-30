package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoJavierMartin/devboard/internal/domain"
	generated "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	var taskID, authorID pgtype.UUID

	if err := taskID.Scan(comment.TaskID); err != nil {
		return fmt.Errorf("invalid task_id: %w", domain.ErrInvalidInput)
	}

	if err := authorID.Scan(comment.AuthorID); err != nil {
		return fmt.Errorf("invalid author_id: %w", domain.ErrInvalidInput)
	}

	result, err := r.queries.CreateComment(ctx, generated.CreateCommentParams{
		TaskID:   taskID,
		AuthorID: authorID,
		Content:  comment.Content,
	})

	if err != nil {
		return fmt.Errorf("inserting comment into DB: %w", err)
	}

	comment.ID = result.ID.String()
	comment.CreatedAt = result.CreatedAt.Time
	comment.UpdatedAt = result.UpdatedAt.Time

	return nil

}

func (r *CommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id: %s: %w", id, domain.ErrNotFound)
	}

	comment, err := r.queries.GetCommentByID(ctx, pgID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("comment %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("retrieving comment %s: %w")
	}

	return toDomainComment(comment), nil
}

func (r *CommentRepository) ListByTask(ctx context.Context, taskID string) ([]*domain.Comment, error) {
	var pgTaskID pgtype.UUID

	if err := pgTaskID.Scan(taskID); err != nil {
		return nil, fmt.Errorf("invalid task_id: %w", domain.ErrInvalidInput)
	}

}

func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	var pgID pgtype.UUID

	if err := pgID.Scan(id); err != nil {
		return fmt.Errorf("invalid id: %w", domain.ErrInvalidInput)
	}

}

func toDomainComment(c generated.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        c.ID.String(),
		TaskID:    c.TaskID.String(),
		AuthorID:  c.AuthorID.String(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}

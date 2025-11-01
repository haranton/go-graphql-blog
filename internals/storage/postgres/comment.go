package postgres

import (
	"context"
	"fmt"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/lib/pq"
)

func (st *PostgresStorage) Comment(ctx context.Context, id int) (*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1;
	`
	var comment models.Comment
	err := st.db.GetContext(ctx, &comment, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

func (st *PostgresStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	query := `
        INSERT INTO comments (post_id, parent_id, user_id, content, created_at, updated_at)
        VALUES (:post_id, :parent_id, :user_id, :content, NOW(), NOW())
        RETURNING id
    `
	if err := execInsertReturningID(ctx, st.db, query, comment, &comment.ID); err != nil {
		return nil, err
	}
	return comment, nil
}

func (st *PostgresStorage) ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY created_at
		LIMIT $2 OFFSET $3;
	`

	var comments []models.Comment
	err := st.db.SelectContext(ctx, &comments, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	return comments, nil
}

func (st *PostgresStorage) ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error) {
	if len(parentIDs) == 0 {
		return []models.Comment{}, nil
	}

	query := `
		SELECT id, post_id, parent_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE parent_id = ANY($1)
		ORDER BY created_at;
	`
	var comments []models.Comment
	err := st.db.SelectContext(ctx, &comments, query, pq.Array(parentIDs))
	if err != nil {
		return nil, fmt.Errorf("batch load replies: %w", err)
	}
	return comments, nil
}

func (st *PostgresStorage) ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at
		LIMIT $2 OFFSET $3;
	`
	var comments []models.Comment
	err := st.db.SelectContext(ctx, &comments, query, parentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list replies: %w", err)
	}
	return comments, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *PostgresStorage) GetPost(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	err := st.db.GetContext(ctx, &post,
		`SELECT id, title, content, user_id, allow_comments, created_at, updated_at 
		 FROM posts 
		 WHERE id = $1`, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (st *PostgresStorage) Posts(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT id, title, content, user_id, allow_comments, created_at, updated_at 
			  FROM posts 
			  ORDER BY created_at DESC`
	err := st.db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	return posts, nil
}

func (st *PostgresStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	query := `
        INSERT INTO posts (title, content, user_id, allow_comments, created_at, updated_at)
        VALUES (:title, :content, :user_id, :allow_comments, NOW(), NOW())
        RETURNING id
    `
	if err := execInsertReturningID(ctx, st.db, query, post, &post.ID); err != nil {
		return nil, err
	}
	return post, nil
}

func (st *PostgresStorage) SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error {
	var ownerID int
	err := st.db.GetContext(ctx, &ownerID, `SELECT user_id FROM posts WHERE id = $1`, postID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("post not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get post owner: %w", err)
	}

	if ownerID != userID {
		return fmt.Errorf("forbidden: not post owner")
	}

	_, err = st.db.ExecContext(ctx,
		`UPDATE posts SET allow_comments = $1, updated_at = NOW() WHERE id = $2`,
		allow, postID)
	if err != nil {
		return fmt.Errorf("failed to update allow_comments: %w", err)
	}
	return nil
}

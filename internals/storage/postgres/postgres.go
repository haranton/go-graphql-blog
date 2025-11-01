package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (st *PostgresStorage) GetPost(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	err := st.db.GetContext(ctx, &post,
		`SELECT id, title, content, user_id, allow_comments, created_at, updated_at 
		 FROM posts 
		 WHERE id = $1`, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // пост не найден — не ошибка
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

// CreateComment создаёт комментарий и возвращает его с ID
func (st *PostgresStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	if comment.PostID == 0 {
		return nil, fmt.Errorf("post_id is required")
	}
	if len([]rune(comment.Content)) > 2000 {
		return nil, fmt.Errorf("comment content exceeds 2000 characters")
	}

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

// ListRepliesBatch — батчевая выборка для DataLoader
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

// ListReplies — одиночная выборка (удобно для тестов и сервисного слоя)
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

func (st *PostgresStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (login, password)
		VALUES (:login, :password)
		RETURNING id
	`
	if err := execInsertReturningID(ctx, st.db, query, user, &user.ID); err != nil {
		return nil, err
	}
	return user, nil
}

func (st *PostgresStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := st.db.GetContext(ctx, &user, `SELECT id, login, password FROM users WHERE login = $1`, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

func (st *PostgresStorage) Close() error {
	return st.db.Close()
}

func execInsertReturningID(ctx context.Context, db *sqlx.DB, query string, obj interface{}, idPtr *int) error {
	rows, err := db.NamedQueryContext(ctx, query, obj)
	if err != nil {
		return fmt.Errorf("exec insert with returning id failed: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(idPtr); err != nil {
			return fmt.Errorf("scan returned id failed: %w", err)
		}
		return nil
	}
	return fmt.Errorf("no id returned after insert")
}

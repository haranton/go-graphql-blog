package postgres

import (
	"context"
	"fmt"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/jmoiron/sqlx"
)

// type Storage interface {
// 	Posts(ctx context.Context) ([]models.Post, error)
// 	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
// 	PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error)
// 	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
// 	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
// 	UserByLogin(ctx context.Context, login string) (*models.User, error)
// 	ListComments(ctx context.Context, postID int, parentID *int, limit, offset int) ([]models.Comment, error)
// 	SetPostAllowComments(ctx context.Context, postID int, userID int, allow bool) error
// 	Close() error
// }

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (st *PostgresStorage) Posts(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	query := `
		SELECT id, title, content, user_id, allow_comments, created_at, updated_at
		FROM posts
		ORDER BY created_at
	`

	if err := st.db.SelectContext(ctx, &posts, query); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return posts, nil
}

func (st *PostgresStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	query := `
        INSERT INTO comments (post_id, parent_id, user_id, content)
		VALUES (:post_id, :parent_id, :user_id, :content)
		RETURNING id
    `

	rows, err := st.db.NamedQueryContext(ctx, query, comment)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&comment.ID); err != nil {
			return nil, fmt.Errorf("failed to scan returned id: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no id returned after insert")
	}

	return comment, nil
}

func (st *PostgresStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	query := `
        INSERT INTO posts (title, content, user_id, allow_comments, created_at, updated_at)
        VALUES (:title, :content, :user_id, :allow_comments, NOW(), NOW())
        RETURNING id
    `

	rows, err := st.db.NamedQueryContext(ctx, query, post)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&post.ID); err != nil {
			return nil, fmt.Errorf("failed to scan returned id: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no id returned after insert")
	}

	return post, nil
}

func (st *PostgresStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (login, password)
		VALUES (:login, :password)
		RETURNING id
	`

	rows, err := st.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return nil, fmt.Errorf("failed to scan returned id: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no id returned after insert")
	}

	return user, nil
}

func (st *PostgresStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user *models.User

	query := `
		SELECT id, login, password
		FROM users
		WHERE login = $1
	`

	if err := st.db.SelectContext(ctx, &user, query, login); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return user, nil
}

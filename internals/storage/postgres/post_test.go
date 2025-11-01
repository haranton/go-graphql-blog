package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage/postgres"
)

func TestCreatePost_Success(t *testing.T) {
	_, mock, store := SetupMock(t)

	post := &models.Post{
		Title:         "Test title",
		Content:       "Test content",
		UserID:        1,
		AllowComments: true,
	}

	postNewID := 10

	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO posts (title, content, user_id, allow_comments, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
        RETURNING id
    `)).
		WithArgs(post.Title, post.Content, post.UserID, post.AllowComments).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(postNewID))

	result, err := store.CreatePost(context.Background(), post)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 10, result.ID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePost_DataNotFilled(t *testing.T) {
	_, mock, store := SetupMock(t)

	post := &models.Post{
		Title:   "",
		Content: "",
		UserID:  1,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO posts (title, content, user_id, allow_comments, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
        RETURNING id
    `)).
		WithArgs(post.Title, post.Content, post.UserID, post.AllowComments).WillReturnError(fmt.Errorf("insert failed"))

	result, err := store.CreatePost(context.Background(), post)
	require.Error(t, err)
	require.Nil(t, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSetPostAllowComments_Success(t *testing.T) {
	_, mock, store := SetupMock(t)

	postID := 1
	userID := 5
	allow := false

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id FROM posts WHERE id = $1`)).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE posts SET allow_comments = $1, updated_at = NOW() WHERE id = $2`)).
		WithArgs(allow, postID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := store.SetPostAllowComments(context.Background(), postID, userID, allow)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSetPostAllowComments_PostNotFound(t *testing.T) {
	_, mock, store := SetupMock(t)

	postID := 1
	userID := 5
	allow := false

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id FROM posts WHERE id = $1`)).
		WithArgs(postID).WillReturnError(sql.ErrNoRows)

	err := store.SetPostAllowComments(context.Background(), postID, userID, allow)
	require.Error(t, err)
	require.Contains(t, err.Error(), "post not found")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSetPostAllowComments_Forbidden(t *testing.T) {
	_, mock, store := SetupMock(t)

	postID := 1
	ownerID := 10
	userID := 5
	allow := false

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id FROM posts WHERE id = $1`)).
		WithArgs(postID).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(ownerID))

	err := store.SetPostAllowComments(context.Background(), postID, userID, allow)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden: not post owner")
	require.NoError(t, mock.ExpectationsWereMet())
}

func SetupMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *postgres.PostgresStorage) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := postgres.NewPostgresStorage(sqlxDB)
	return sqlxDB, mock, store
}

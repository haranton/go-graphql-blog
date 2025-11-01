package postgres

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateComment_Success(t *testing.T) {
	_, mock, store := SetupMock(t)

	comment := &models.Comment{
		PostID:   1,
		ParentID: nil,
		UserID:   1,
		Content:  "Test content",
	}
	commentNewID := 10

	mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO comments (post_id, parent_id, user_id, content, created_at, updated_at)
    VALUES (?, ?, ?, ?, NOW(), NOW())
    RETURNING id
`)).
		WithArgs(comment.PostID, comment.ParentID, comment.UserID, comment.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(commentNewID))

	result, err := store.CreateComment(context.Background(), comment)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, commentNewID, result.ID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateComment_dataNotFilled(t *testing.T) {
	_, mock, store := SetupMock(t)

	comment := &models.Comment{
		ParentID: nil,
		UserID:   1,
		Content:  "Test content",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
    INSERT INTO comments (post_id, parent_id, user_id, content, created_at, updated_at)
    VALUES (?, ?, ?, ?, NOW(), NOW())
    RETURNING id
`)).
		WithArgs(comment.PostID, comment.ParentID, comment.UserID, comment.Content).
		WillReturnError(sqlmock.ErrCancelled)

	result, err := store.CreateComment(context.Background(), comment)
	require.Error(t, err)
	require.Nil(t, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func SetupMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *PostgresStorage) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	store := NewPostgresStorage(sqlxDB)
	return sqlxDB, mock, store
}

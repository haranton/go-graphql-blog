package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/haranton/go-graphql-blog/internals/auth"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentService_AddComment(t *testing.T) {
	t.Run("successful creation with parent", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		parentIDStr := "5"
		expectedComment := &models.Comment{
			ID:       1,
			PostID:   10,
			ParentID: func() *int { id := 5; return &id }(),
			UserID:   1,
			Content:  "Test comment",
		}

		mockStore.CreateCommentFunc = func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
			assert.Equal(t, 10, comment.PostID)
			require.NotNil(t, comment.ParentID)
			assert.Equal(t, 5, *comment.ParentID)
			assert.Equal(t, 1, comment.UserID)
			assert.Equal(t, "Test comment", comment.Content)
			return expectedComment, nil
		}

		result, err := service.AddComment(ctx, "10", &parentIDStr, "Test comment")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "1", result.ID)
	})

	t.Run("successful creation without parent", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		mockStore.CreateCommentFunc = func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
			assert.Nil(t, comment.ParentID)
			return &models.Comment{ID: 1}, nil
		}

		result, err := service.AddComment(ctx, "10", nil, "Test comment")
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("empty content", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		result, err := service.AddComment(ctx, "10", nil, "")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("content too long", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		longContent := string(make([]rune, 2001))
		result, err := service.AddComment(ctx, "10", nil, longContent)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid post ID", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		result, err := service.AddComment(ctx, "invalid", nil, "Test comment")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid parent ID", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		invalidParentID := "invalid"
		result, err := service.AddComment(ctx, "10", &invalidParentID, "Test comment")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		mockStore.CreateCommentFunc = func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
			return nil, errors.New("storage error")
		}

		result, err := service.AddComment(ctx, "10", nil, "Test comment")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCommentService_GetComments(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		log := slog.New(slog.NewTextHandler(os.Stdout, nil))
		mockStore := &mocks.MockStorage{}
		service := NewCommentService(mockStore, log)

		expectedComments := []models.Comment{
			{ID: 1, Content: "Comment 1"},
			{ID: 2, Content: "Comment 2"},
		}

		mockStore.ListCommentsFunc = func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
			assert.Equal(t, 10, postID)
			assert.Equal(t, 20, limit)
			assert.Equal(t, 5, offset)
			return expectedComments, nil
		}

		result, err := service.GetComments(context.Background(), "10", 20, 5)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ID)
	})

	t.Run("limit validation - default to 10", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)

		mockStore.ListCommentsFunc = func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
			assert.Equal(t, 10, limit) // default limit
			return []models.Comment{}, nil
		}

		_, err := service.GetComments(context.Background(), "10", 0, 0)
		require.NoError(t, err)
	})

	t.Run("limit validation - max 100", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)
		mockStore.ListCommentsFunc = func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
			assert.Equal(t, maxCommentsPageSize, limit)
			return []models.Comment{}, nil
		}

		_, err := service.GetComments(context.Background(), "10", 200, 0)
		require.NoError(t, err)
	})

	t.Run("offset validation - negative to 0", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)

		mockStore.ListCommentsFunc = func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
			assert.Equal(t, 0, offset)
			return []models.Comment{}, nil
		}

		_, err := service.GetComments(context.Background(), "10", 10, -5)
		require.NoError(t, err)
	})

	t.Run("invalid post ID", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := service(mockStore)

		result, err := service.GetComments(context.Background(), "invalid", 10, 0)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("storage error", func(t *testing.T) {

		mockStore := &mocks.MockStorage{}
		service := service(mockStore)

		mockStore.ListCommentsFunc = func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
			return nil, errors.New("storage error")
		}

		result, err := service.GetComments(context.Background(), "10", 10, 0)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "could not list comments")
	})
}

func service(mockStore *mocks.MockStorage) *commentService {
	return &commentService{
		store:   mockStore,
		slogger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

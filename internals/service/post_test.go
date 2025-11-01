package service

import (
	"context"
	"errors"
	"testing"

	"github.com/haranton/go-graphql-blog/internals/auth"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostService_CreatePost(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)

		userID := 1
		ctx := auth.WithContext(context.Background(), &models.User{ID: userID})

		expectedPost := &models.Post{
			ID:            1,
			Title:         "Test Post",
			Content:       "Content",
			UserID:        userID,
			AllowComments: true,
		}

		mockStore.CreatePostFunc = func(ctx context.Context, post *models.Post) (*models.Post, error) {
			assert.Equal(t, "Test Post", post.Title)
			assert.Equal(t, "Content", post.Content)
			assert.Equal(t, userID, post.UserID)
			assert.True(t, post.AllowComments)
			return expectedPost, nil
		}

		result, err := service.CreatePost(ctx, "Test Post", "Content")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "Test Post", result.Title)
	})

	t.Run("empty title", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		result, err := service.CreatePost(ctx, "", "Content")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty content", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		result, err := service.CreatePost(ctx, "Title", "")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		mockStore.CreatePostFunc = func(ctx context.Context, post *models.Post) (*models.Post, error) {
			return nil, errors.New("storage error")
		}

		result, err := service.CreatePost(ctx, "Title", "Content")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestPostService_PostWithComments(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)

		expectedPost := &models.Post{
			ID:    5,
			Title: "Test Post",
		}

		mockStore.GetPostFunc = func(ctx context.Context, id int) (*models.Post, error) {
			assert.Equal(t, 5, id)
			return expectedPost, nil
		}

		result, err := service.Post(context.Background(), "5")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "5", result.ID)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)

		result, err := service.Post(context.Background(), "invalid")
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("post not found", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)

		mockStore.GetPostFunc = func(ctx context.Context, id int) (*models.Post, error) {
			return nil, nil
		}

		result, err := service.Post(context.Background(), "999")
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestPostService_DisallowComments(t *testing.T) {
	t.Run("successful disallow", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		mockStore.SetPostAllowCommentsFunc = func(ctx context.Context, postID, userID int, allow bool) error {
			assert.Equal(t, 10, postID)
			assert.Equal(t, 1, userID)
			assert.False(t, allow)
			return nil
		}

		result, err := service.DisallowComments(ctx, "10")
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("invalid post ID", func(t *testing.T) {
		mockStore := &mocks.MockStorage{}
		service := NewPostService(mockStore)
		ctx := auth.WithContext(context.Background(), &models.User{ID: 1})

		result, err := service.DisallowComments(ctx, "invalid")
		require.Error(t, err)
		assert.False(t, result)
	})
}

package memory

import (
	"context"
	"testing"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	st := NewMemoryStorage()
	ctx := context.Background()

	post := models.Post{
		Title:   "Test Post",
		Content: "Some content",
		UserID:  42,
	}

	created, err := st.CreatePost(ctx, &post)
	require.NoError(t, err)

	require.Equal(t, 1, created.ID)
	require.Equal(t, "Test Post", created.Title)
	require.Equal(t, "Some content", created.Content)
	require.Equal(t, 42, created.UserID)
	require.NotZero(t, created.CreatedAt)
	require.NotZero(t, created.UpdatedAt)

	require.Len(t, st.posts, 1)
	require.Contains(t, st.postsByID, created.ID)
}

func TestPosts_ReturnsAll(t *testing.T) {
	st := NewMemoryStorage()
	ctx := context.Background()

	st.CreatePost(ctx, &models.Post{Title: "One", Content: "A", UserID: 1})
	st.CreatePost(ctx, &models.Post{Title: "Two", Content: "B", UserID: 2})

	posts, err := st.Posts(ctx)
	require.NoError(t, err)
	require.Len(t, posts, 2)
	require.Equal(t, "One", posts[0].Title)
	require.Equal(t, "Two", posts[1].Title)

	// Проверяем, что вернулись копии
	posts[0].Title = "Changed"
	require.NotEqual(t, posts[0].Title, st.posts[0].Title)
}

func TestGetPost(t *testing.T) {
	st := NewMemoryStorage()
	ctx := context.Background()

	post, _ := st.CreatePost(ctx, &models.Post{
		Title:   "Post",
		Content: "Body",
		UserID:  1,
	})

	got, err := st.GetPost(ctx, post.ID)
	require.NoError(t, err)
	require.Equal(t, post.Title, got.Title)
	require.Equal(t, post.Content, got.Content)
	require.NotSame(t, post, got)

	_, err = st.GetPost(ctx, 999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "post not found")
}

func TestSetPostAllowComments(t *testing.T) {
	st := NewMemoryStorage()
	ctx := context.Background()

	post, _ := st.CreatePost(ctx, &models.Post{
		Title:         "Sample",
		Content:       "Text",
		UserID:        10,
		AllowComments: true,
	})

	err := st.SetPostAllowComments(ctx, post.ID, 10, false)
	require.NoError(t, err)
	require.False(t, st.postsByID[post.ID].AllowComments)

	err = st.SetPostAllowComments(ctx, post.ID, 99, true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden")

	err = st.SetPostAllowComments(ctx, 999, 10, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "post not found")
}

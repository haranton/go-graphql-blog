package memory

import (
	"context"
	"testing"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/stretchr/testify/require"
)

func TestCreateComment(t *testing.T) {
	storage := NewMemoryStorage()

	commentName := "test"
	comment := models.Comment{
		PostID:  1,
		Content: commentName,
		UserID:  2,
	}

	created, err := storage.CreateComment(context.TODO(), &comment)
	require.NoError(t, err)
	require.Equal(t, created.Content, commentName)
	require.NotZero(t, created.CreatedAt)
	require.Equal(t, 1, created.ID)
}

func TestComment(t *testing.T) {
	storage := NewMemoryStorage()

	idPost := 1

	commentName := "test"
	comment := models.Comment{
		PostID:  idPost,
		Content: commentName,
		UserID:  2,
	}

	storage.CreateComment(context.TODO(), &comment)

	commentCreated, err := storage.Comment(context.TODO(), idPost)
	require.NoError(t, err)
	require.Equal(t, 1, commentCreated.ID)

}

func TestListComments_TopLevelOnly(t *testing.T) {
	st := NewMemoryStorage()
	ctx := context.Background()

	parent, _ := st.CreateComment(ctx, &models.Comment{PostID: 1, Content: "Parent"})
	reply, _ := st.CreateComment(ctx, &models.Comment{PostID: 1, ParentID: &parent.ID, Content: "Child"})

	list, err := st.ListComments(ctx, 1, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Parent", list[0].Content)

	require.NotContains(t, list, reply.Content)
}

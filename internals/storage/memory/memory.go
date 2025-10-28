package memory

import (
	"context"
	"sync"

	"github.com/haranton/go-graphql-blog/internals/models"
)

// type Storage interface {
// 	Posts(ctx context.Context) ([]models.Post, error)
// 	CreatePost(ctx context.Context, post *models.Post) (models.Post, error)
// 	PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error) //todo
// 	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
// }

type MemoryStorage struct {
	posts         []*models.Post
	comments      []*models.Comment
	mu            sync.Mutex
	nextPostID    int
	nextCommentID int
}

func NewMemoryStorage() *MemoryStorage {

	return &MemoryStorage{
		posts:    []*models.Post{},
		comments: []*models.Comment{},
	}
}

func (st *MemoryStorage) Post(ctx context.Context) ([]models.Post, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	resultPosts := make([]models.Post, len(st.posts))
	for i, postPointer := range st.posts {
		resultPosts[i] = *postPointer
	}

	return resultPosts, nil
}

func (st *MemoryStorage) CreatePost(ctx context.Context, post *models.Post) (models.Post, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	post.ID = st.nextPostID + 1
	st.posts = append(st.posts, post)
	st.nextPostID = st.nextPostID + 1

	return *post, nil
}

func (st *MemoryStorage) PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error) {
	st.mu.Lock()

}

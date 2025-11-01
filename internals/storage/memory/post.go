package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *MemoryStorage) Posts(ctx context.Context) ([]models.Post, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	result := make([]models.Post, len(st.posts))
	for i, postPtr := range st.posts {
		result[i] = *postPtr
	}
	return result, nil
}

func (st *MemoryStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextPostID++
	post.ID = st.nextPostID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	copy := *post
	st.posts = append(st.posts, &copy)

	return &copy, nil
}

func (st *MemoryStorage) GetPost(ctx context.Context, idPost int) (*models.Post, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var post *models.Post
	for _, p := range st.posts {
		if p.ID == idPost {
			post = p
			break
		}
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

func (st *MemoryStorage) SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, p := range st.posts {
		if p.ID == postID {
			if p.UserID != userID {
				return fmt.Errorf("forbidden: not post owner")
			}
			p.AllowComments = allow
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("post not found")
}

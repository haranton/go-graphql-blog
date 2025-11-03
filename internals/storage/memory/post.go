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
	st.postsByID[copy.ID] = &copy

	return &copy, nil
}

func (st *MemoryStorage) GetPost(ctx context.Context, idPost int) (*models.Post, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if p, ok := st.postsByID[idPost]; ok {
		copy := *p
		return &copy, nil
	}
	return nil, fmt.Errorf("post not found")
}

func (st *MemoryStorage) SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	p, ok := st.postsByID[postID]
	if !ok {
		return fmt.Errorf("post not found")
	}
	if p.UserID != userID {
		return fmt.Errorf("forbidden: not post owner")
	}
	p.AllowComments = allow
	p.UpdatedAt = time.Now()
	return nil
}

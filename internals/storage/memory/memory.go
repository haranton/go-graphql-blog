package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type MemoryStorage struct {
	posts         []*models.Post
	comments      []*models.Comment
	users         []*models.User
	mu            sync.Mutex
	nextPostID    int
	nextCommentID int
	nextUserID    int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    []*models.Post{},
		comments: []*models.Comment{},
		users:    []*models.User{},
	}
}

func (st *MemoryStorage) Posts(ctx context.Context) ([]models.Post, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

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
	st.mu.Lock()
	defer st.mu.Unlock()

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

func (st *MemoryStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextCommentID++
	comment.ID = st.nextCommentID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	copy := *comment
	st.comments = append(st.comments, &copy)

	return &copy, nil
}

// ListComments — верхний уровень комментариев (parent_id IS NULL)
func (st *MemoryStorage) ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.PostID == postID && c.ParentID == nil {
			filtered = append(filtered, *c)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	if offset >= len(filtered) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

// ListReplies — одиночный parent_id (для сервисного слоя)
func (st *MemoryStorage) ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.ParentID != nil && *c.ParentID == parentID {
			filtered = append(filtered, *c)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	if offset >= len(filtered) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

// ListRepliesBatch — для DataLoader
func (st *MemoryStorage) ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	ids := make(map[int]struct{}, len(parentIDs))
	for _, id := range parentIDs {
		ids[id] = struct{}{}
	}

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.ParentID != nil {
			if _, ok := ids[*c.ParentID]; ok {
				filtered = append(filtered, *c)
			}
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	return filtered, nil
}

func (st *MemoryStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextUserID++
	user.ID = st.nextUserID
	st.users = append(st.users, user)

	return user, nil
}

func (st *MemoryStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, u := range st.users {
		if u.Login == login {
			return u, nil
		}
	}
	return nil, nil
}

func (st *MemoryStorage) Close() error { return nil }

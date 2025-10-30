package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/haranton/go-graphql-blog/internals/models"
)

// type Storage interface {
// 	Posts(ctx context.Context) ([]models.Post, error)
// 	CreatePost(ctx context.Context, post *models.Post) (models.Post, error)
// 	PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error) //todo
// 	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
// }

//todo Оптимизация поиска

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

	resultPosts := make([]models.Post, len(st.posts))
	for i, postPointer := range st.posts {
		resultPosts[i] = *postPointer
	}

	return resultPosts, nil
}

func (st *MemoryStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	post.ID = st.nextPostID + 1
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	st.nextPostID++

	stored := *post
	st.posts = append(st.posts, &stored)

	result := stored

	return &result, nil
}

func (st *MemoryStorage) PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error) {
	st.mu.Lock()

	isFindPost := false
	var resultPostWithComments models.PostWithComments
	for _, post := range st.posts {
		if post.ID == idPost {
			isFindPost = true
			resultPostWithComments.AllowComments = post.AllowComments
			resultPostWithComments.UserID = post.UserID
			resultPostWithComments.Content = post.Content
			resultPostWithComments.CreatedAt = post.CreatedAt
			resultPostWithComments.ID = post.ID
			resultPostWithComments.Title = post.Title
			resultPostWithComments.UpdatedAt = post.UpdatedAt
		}
	}

	if !isFindPost {
		st.mu.Unlock()
		return nil, nil
	}

	st.mu.Unlock()
	return &resultPostWithComments, nil //todo убраны комментарии
}

func (st *MemoryStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	comment.ID = st.nextCommentID + 1
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	st.nextCommentID++

	stored := *comment
	st.comments = append(st.comments, &stored)

	result := stored

	return &result, nil
}

func (st *MemoryStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	user.ID = st.nextUserID + 1
	st.nextUserID++
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

func (st *MemoryStorage) ListComments(ctx context.Context, postID int, parentID *int, limit, offset int) ([]models.Comment, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	var filtered []models.Comment
	for _, c := range st.comments {
		if c.PostID != postID {
			continue
		}

		// Проверяем соответствие parentID
		matchesParent := (parentID == nil && c.ParentID == nil) ||
			(parentID != nil && c.ParentID != nil && *c.ParentID == *parentID)

		if matchesParent {
			filtered = append(filtered, *c)
		}
	}

	// Сортируем по времени создания (если есть поле CreatedAt)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	// Пагинация
	if offset >= len(filtered) {
		return []models.Comment{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (st *MemoryStorage) SetPostAllowComments(ctx context.Context, postID int, userID int, allow bool) error {
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

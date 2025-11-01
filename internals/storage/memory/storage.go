package memory

import (
	"sync"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type MemoryStorage struct {
	posts         []*models.Post
	comments      []*models.Comment
	users         []*models.User
	mu            sync.RWMutex
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

func (st *MemoryStorage) Close() error { return nil }

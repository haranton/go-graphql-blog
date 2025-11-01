package memory

import (
	"sync"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type MemoryStorage struct {
	posts    []*models.Post
	comments []*models.Comment
	users    []*models.User

	postsByID   map[int]*models.Post
	commentByID map[int]*models.Comment
	userBylogin map[string]*models.User

	commentsByPostID  map[int][]*models.Comment
	repliesByParentID map[int][]*models.Comment

	mu            sync.RWMutex
	nextPostID    int
	nextCommentID int
	nextUserID    int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		postsByID:         make(map[int]*models.Post),
		commentByID:       make(map[int]*models.Comment),
		userBylogin:       make(map[string]*models.User),
		commentsByPostID:  make(map[int][]*models.Comment),
		repliesByParentID: make(map[int][]*models.Comment),
	}
}

func (st *MemoryStorage) Close() error { return nil }

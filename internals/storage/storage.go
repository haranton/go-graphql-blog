package storage

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type Storage interface {
	// Posts
	Posts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error

	// Comments
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error)
	ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error)
	ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error)

	// Users
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UserByLogin(ctx context.Context, login string) (*models.User, error)

	Close() error
}

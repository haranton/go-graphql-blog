package storage

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type DB interface {
	Close() error
}

type PostStorage interface {
	Posts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error
}

type CommentStorage interface {
	Comment(ctx context.Context, id int) (*models.Comment, error)
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error)
	ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error)
	ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error)
}

type UserStorage interface {
	UserByLogin(ctx context.Context, login string) (*models.User, error)
}

type Storage interface {
	PostStorage
	CommentStorage
	UserStorage
	DB
}

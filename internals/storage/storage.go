package storage

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type Storage interface {
	Posts(ctx context.Context) ([]models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	// PostWithComments(ctx context.Context, idPost int) (*models.PostWithComments, error) //todo
	// CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
}

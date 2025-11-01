package mocks

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

// MockStorage is a mock implementation of storage.Storage
type MockStorage struct {
	// Posts
	PostsFunc                func(ctx context.Context) ([]models.Post, error)
	GetPostFunc              func(ctx context.Context, id int) (*models.Post, error)
	CreatePostFunc           func(ctx context.Context, post *models.Post) (*models.Post, error)
	SetPostAllowCommentsFunc func(ctx context.Context, postID, userID int, allow bool) error

	// Comments
	CommentFunc          func(ctx context.Context, id int) (*models.Comment, error)
	CreateCommentFunc    func(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	ListCommentsFunc     func(ctx context.Context, postID, limit, offset int) ([]models.Comment, error)
	ListRepliesBatchFunc func(ctx context.Context, parentIDs []int) ([]models.Comment, error)
	ListRepliesFunc      func(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error)

	// Users
	CreateUserFunc  func(ctx context.Context, user *models.User) (*models.User, error)
	UserByLoginFunc func(ctx context.Context, login string) (*models.User, error)

	CloseFunc func() error
}

var _ storage.Storage = (*MockStorage)(nil)

func (m *MockStorage) Posts(ctx context.Context) ([]models.Post, error) {
	if m.PostsFunc != nil {
		return m.PostsFunc(ctx)
	}
	return nil, nil
}

func (m *MockStorage) GetPost(ctx context.Context, id int) (*models.Post, error) {
	if m.GetPostFunc != nil {
		return m.GetPostFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	if m.CreatePostFunc != nil {
		return m.CreatePostFunc(ctx, post)
	}
	return nil, nil
}

func (m *MockStorage) SetPostAllowComments(ctx context.Context, postID, userID int, allow bool) error {
	if m.SetPostAllowCommentsFunc != nil {
		return m.SetPostAllowCommentsFunc(ctx, postID, userID, allow)
	}
	return nil
}

func (m *MockStorage) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	if m.CreateCommentFunc != nil {
		return m.CreateCommentFunc(ctx, comment)
	}
	return nil, nil
}

func (m *MockStorage) ListComments(ctx context.Context, postID, limit, offset int) ([]models.Comment, error) {
	if m.ListCommentsFunc != nil {
		return m.ListCommentsFunc(ctx, postID, limit, offset)
	}
	return nil, nil
}

func (m *MockStorage) ListRepliesBatch(ctx context.Context, parentIDs []int) ([]models.Comment, error) {
	if m.ListRepliesBatchFunc != nil {
		return m.ListRepliesBatchFunc(ctx, parentIDs)
	}
	return nil, nil
}

func (m *MockStorage) ListReplies(ctx context.Context, parentID, limit, offset int) ([]models.Comment, error) {
	if m.ListRepliesFunc != nil {
		return m.ListRepliesFunc(ctx, parentID, limit, offset)
	}
	return nil, nil
}

func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil, nil
}

func (m *MockStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	if m.UserByLoginFunc != nil {
		return m.UserByLoginFunc(ctx, login)
	}
	return nil, nil
}

func (m *MockStorage) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockStorage) Comment(ctx context.Context, id int) (*models.Comment, error) {
	if m.CommentFunc != nil {
		return m.CommentFunc(ctx, id)
	}
	return nil, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/haranton/go-graphql-blog/graph/model"
	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/auth"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

const MaxCommentsPageSize = 100

type commentService struct {
	store storage.Storage
}

func NewCommentService(store storage.Storage) *commentService {
	return &commentService{store: store}
}

func (s *commentService) AddComment(ctx context.Context, postIDStr string, parentIDStr *string, content string) (*gqlmodel.Comment, error) {
	if content == "" {
		return nil, errors.New("content must be provided")
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return nil, err
	}

	contentForCheckLen := []rune(content)
	if len(contentForCheckLen) > 2000 {
		return nil, errors.New("content length must be less than 2000 characters")
	}

	var parentID *int
	if parentIDStr != nil && *parentIDStr != "" {
		pid, err2 := strconv.Atoi(*parentIDStr)
		if err2 != nil {
			return nil, err2
		}
		parentID = &pid
	}

	user := auth.ForContext(ctx)

	domainComment := &models.Comment{
		PostID:   postID,
		ParentID: parentID,
		UserID:   user.ID,
		Content:  content,
	}

	created, err := s.store.CreateComment(ctx, domainComment)
	if err != nil {
		return nil, err
	}

	return mapper.MapCommentDomainToGraph(created), nil
}

func (s *commentService) GetComments(ctx context.Context, postIDStr string, limit int, offset int) ([]*model.Comment, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > MaxCommentsPageSize {
		limit = MaxCommentsPageSize
	}
	if offset < 0 {
		offset = 0
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid postID: %w", err)
	}

	comments, err := s.store.ListComments(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list comments: %w", err)
	}

	result := make([]*model.Comment, len(comments))
	for i, c := range comments {
		result[i] = mapper.MapCommentDomainToGraph(&c)
	}
	return result, nil
}

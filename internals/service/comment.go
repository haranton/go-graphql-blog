package service

import (
	"context"
	"errors"
	"strconv"

	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/auth"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

type commentService struct {
	store storage.Storage
}

func NewCommentService(store storage.Storage) *commentService {
	return &commentService{store: store}
}

func (s *commentService) AddComment(ctx context.Context, postIDStr string, parentIDStr *string, content string) (*gqlmodel.Comment, error) {
	if content == "" { // возможно излишняя проверка todo
		return nil, errors.New("content must be provided")
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return nil, err
	}

	var parentID *int
	if parentIDStr != nil {
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

	return mapper.MapCommentDomainToGraph(created, []models.Comment{}), nil
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/haranton/go-graphql-blog/graph/model"
	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/auth"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

const maxCommentsPageSize = 100
const defaultOffset = 0
const defaultLimit = 10

type commentService struct {
	store   storage.Storage
	slogger *slog.Logger
}

func NewCommentService(store storage.Storage, logger *slog.Logger) *commentService {
	return &commentService{store: store, slogger: logger}
}

func (s *commentService) AddComment(ctx context.Context, postIDStr string, parentIDStr *string, content string) (*gqlmodel.Comment, error) {

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return nil, err
	}

	var parentID *int
	if parentIDStr != nil && *parentIDStr != "" {
		pid, err2 := strconv.Atoi(*parentIDStr)
		if err2 != nil {
			return nil, err2
		}
		parentID = &pid
	}

	//Проверяем что post сущеcтвует
	post, err := s.store.GetPost(ctx, postID)
	if err != nil {
		s.slogger.Error("failed to get post", "postID", postID, "error", err)
		return nil, err
	}

	if post == nil {
		s.slogger.Error("post not found", "postID", postID)
		return nil, fmt.Errorf("post with ID %d not found", postID)
	}

	if !post.AllowComments {
		s.slogger.Error("comments are disabled", "postID", postID)
		return nil, fmt.Errorf("comments are disabled for post with ID %d", postID)
	}

	//Проверяем что комментарий с parentId сущевствует
	if parentID != nil {
		parentComment, err := s.store.Comment(ctx, *parentID)
		if err != nil {
			s.slogger.Error("failed to get parent comment", "parentID", *parentID, "error", err)
			return nil, err
		}
		if parentComment == nil {
			s.slogger.Error("parent comment not found", "parentID", *parentID)
			return nil, fmt.Errorf("parent comment with ID %d not found", *parentID)
		}
		if parentComment.PostID != postID {
			s.slogger.Error("parent comment does not belong to post", "parentID", *parentID, "postID", postID)
			return nil, fmt.Errorf("parent comment with ID %d does not belong to post with ID %d", *parentID, postID)
		}
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
		limit = defaultLimit
	}
	if limit > maxCommentsPageSize {
		limit = maxCommentsPageSize
	}
	if offset < 0 {
		offset = defaultOffset
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

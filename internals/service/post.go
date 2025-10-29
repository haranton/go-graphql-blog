package service

import (
	"context"
	"errors"

	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

type PostService struct {
	store storage.Storage
}

func NewPostService(store storage.Storage) *PostService {
	return &PostService{store: store}
}

func (s *PostService) Posts(ctx context.Context) ([]*gqlmodel.Post, error) {
	domainPosts, err := s.store.Posts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*gqlmodel.Post, len(domainPosts))
	for i := range domainPosts {

		result[i] = mapper.MapPostWithCommentsDomainToGraph(
			&models.PostWithComments{
				Post:     domainPosts[i],
				Comments: []models.Comment{},
			},
		)
	}
	return result, nil
}

// Создание нового поста
func (s *PostService) CreatePost(ctx context.Context, title string, content string, author *string) (*gqlmodel.Post, error) {
	if title == "" || content == "" { // todo возможно излично
		return nil, errors.New("title and content must be provided")
	}

	domainPost := &models.Post{
		Title:         title,
		Content:       content,
		Author:        author,
		AllowComments: true,
	}

	created, err := s.store.CreatePost(ctx, domainPost)
	if err != nil {
		return nil, err
	}

	postWithComments := &models.PostWithComments{
		Post:     *created,
		Comments: []models.Comment{},
	}

	return mapper.MapPostWithCommentsDomainToGraph(postWithComments), nil
}

// func (s *PostService) GetPostByID(ctx context.Context, idStr string) (*gqlmodel.Post, error) {
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	domainPW, err := s.store.PostWithComments(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if domainPW == nil {
// 		return nil, nil
// 	}

// 	return mapper.MapPostWithCommentsDomainToGraph(domainPW), nil
// }

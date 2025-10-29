package service

import (
	"context"
	"errors"
	"strconv"

	gqlmodel "github.com/haranton/go-graphql-blog/graph/model"
	"github.com/haranton/go-graphql-blog/internals/mapper"
	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

type postService struct {
	store storage.Storage
}

func NewPostService(store storage.Storage) *postService {
	return &postService{store: store}
}

func (s *postService) Posts(ctx context.Context) ([]*gqlmodel.Post, error) {
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
func (s *postService) CreatePost(ctx context.Context, title string, content string, author *string) (*gqlmodel.Post, error) {
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

func (s *postService) PostWithComments(ctx context.Context, idStr string) (*gqlmodel.Post, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	post, err := s.store.PostWithComments(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil { //todo нужно обрабатывать как ошибку если значение не найдено
		return nil, nil
	}

	return mapper.MapPostWithCommentsDomainToGraph(post), nil
}

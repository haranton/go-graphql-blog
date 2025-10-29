package service

import "github.com/haranton/go-graphql-blog/internals/storage"

type Service struct {
	SrvPost *postService
	SrvComm *commentService
	SrvAuth *AuthService
}

func NewService(storage storage.Storage) *Service {
	srvPost := NewPostService(storage)
	srvComm := NewCommentService(storage)
	srvAuth := NewAuthService(storage)

	return &Service{
		SrvPost: srvPost,
		SrvComm: srvComm,
		SrvAuth: srvAuth,
	}
}

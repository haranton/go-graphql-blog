package service

import (
	"log/slog"

	"github.com/haranton/go-graphql-blog/internals/storage"
)

type Service struct {
	SrvPost *postService
	SrvComm *commentService
	SrvAuth *AuthService
}

func NewService(storage storage.Storage, logger *slog.Logger) *Service {
	srvPost := NewPostService(storage)
	srvComm := NewCommentService(storage, logger)
	srvAuth := NewAuthService(storage)

	return &Service{
		SrvPost: srvPost,
		SrvComm: srvComm,
		SrvAuth: srvAuth,
	}
}

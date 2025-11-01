package graph

import (
	"log/slog"

	"github.com/haranton/go-graphql-blog/internals/service"
	"github.com/haranton/go-graphql-blog/internals/sub"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *service.Service
	Slogger *slog.Logger
	Sub     *sub.CommentPubSub
}

func NewResolver(service *service.Service, slogger *slog.Logger, sub *sub.CommentPubSub) *Resolver {
	return &Resolver{
		Service: service,
		Slogger: slogger,
		Sub:     sub,
	}
}

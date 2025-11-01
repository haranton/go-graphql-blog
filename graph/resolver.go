package graph

import (
	"log/slog"

	"github.com/haranton/go-graphql-blog/internals/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *service.Service
	Slogger *slog.Logger
}

func NewResolver(service *service.Service, slogger *slog.Logger) *Resolver {
	return &Resolver{
		Service: service,
		Slogger: slogger,
	}
}

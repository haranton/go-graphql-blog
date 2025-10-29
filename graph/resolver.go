package graph

import (
	"github.com/haranton/go-graphql-blog/internals/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *service.PostService
}

func NewResolver(service *service.PostService) *Resolver {
	return &Resolver{
		Service: service,
	}
}

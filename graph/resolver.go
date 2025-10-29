package graph

import (
	"github.com/haranton/go-graphql-blog/internals/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *service.Service
}

func NewResolver(service *service.Service) *Resolver {
	return &Resolver{
		Service: service,
	}
}

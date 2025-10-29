package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/haranton/go-graphql-blog/internals/auth"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return next(ctx)
}

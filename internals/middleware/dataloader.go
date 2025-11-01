package middleware

import (
	"context"
	"net/http"

	"github.com/haranton/go-graphql-blog/internals/dataloaders"
	"github.com/haranton/go-graphql-blog/internals/storage"
)

func DataLoaderMiddleware(st storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			loaders := dataloaders.NewCommentLoader(st)

			ctx = context.WithValue(ctx, "commentLoader", loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

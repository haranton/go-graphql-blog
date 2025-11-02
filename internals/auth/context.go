package auth

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type contextKey string

const UserCtxKey contextKey = "user"

func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(UserCtxKey).(*models.User)
	return raw
}

// WithContext adds a user to the context
func WithContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

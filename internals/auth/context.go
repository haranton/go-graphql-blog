package auth

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type contextKey string

const userCtxKey contextKey = "user"

func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(userCtxKey).(*models.User)
	return raw
}

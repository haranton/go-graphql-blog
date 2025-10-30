package auth

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

type contextKey string

const UserCtxKey contextKey = "user"

func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(UserCtxKey).(*models.User)
	if raw == nil {
		return &models.User{
			Login:    "1",
			Password: "1",
			ID:       1, //todo убрать заглушку
		}
	}
	return raw
}

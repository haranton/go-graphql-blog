package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *PostgresStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := st.db.GetContext(ctx, &user, `SELECT id, login, password FROM users WHERE login = $1`, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

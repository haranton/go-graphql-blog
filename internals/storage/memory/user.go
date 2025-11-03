package memory

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *MemoryStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if u, ok := st.userBylogin[login]; ok {
		copy := *u
		return &copy, nil
	}

	for _, u := range st.users {
		if u.Login == login {
			return u, nil
		}
	}
	return nil, nil
}

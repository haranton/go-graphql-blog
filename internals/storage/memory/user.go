package memory

import (
	"context"

	"github.com/haranton/go-graphql-blog/internals/models"
)

func (st *MemoryStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.nextUserID++
	user.ID = st.nextUserID
	st.users = append(st.users, user)

	return user, nil
}

func (st *MemoryStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	for _, u := range st.users {
		if u.Login == login {
			return u, nil
		}
	}
	return nil, nil
}

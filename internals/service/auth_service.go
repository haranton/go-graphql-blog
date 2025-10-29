package service

import (
	"context"
	"errors"

	"github.com/haranton/go-graphql-blog/internals/models"
	"github.com/haranton/go-graphql-blog/internals/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store storage.Storage
}

func NewAuthService(store storage.Storage) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Authenticate(ctx context.Context, login, password string) (*models.User, error) {
	user, err := s.store.UserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user or password uncorrectly") //todo убрать дублирование
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("user or password uncorrectly")
	}

	return user, nil
}

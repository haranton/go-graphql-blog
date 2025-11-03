package service

import (
	"context"
	"errors"
	"fmt"

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
		return nil, errors.New("user or password uncorrectly")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("user or password uncorrectly")
	}

	return user, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, login, password string) (*models.User, error) {

	existing, err := s.store.UserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("login already exists")
	}

	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	userDomain := &models.User{
		Login:    login,
		Password: hashed,
	}
	created, err := s.store.CreateUser(ctx, userDomain)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

package usersvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
	"net/url"
)

type CreateUserInput struct {
	ImageURL  string
	FirstName string
	LastName  string
	Username  string
}

type userService interface {
	CreateUser(ctx context.Context, in CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, in repo.CreateUserInput) error
	GetUser(ctx context.Context, id uuid.UUID) (repo.User, error)
}

type service struct {
	repo userRepository
}

func NewService(repo userRepository) *service {
	return &service{repo}
}

var _ userService = (*service)(nil)

func (s *service) CreateUser(ctx context.Context, in CreateUserInput) (uuid.UUID, error) {
	if in.FirstName == "" {
		return uuid.Nil, errors.New("first name is required")
	}
	if in.Username == "" {
		return uuid.Nil, errors.New("username is required")
	}
	if in.ImageURL == "" {
		return uuid.Nil, errors.New("image url is required")
	}
	u, err := url.ParseRequestURI(in.ImageURL)
	if err != nil || u == nil || u.Scheme == "" || u.Host == "" {
		return uuid.Nil, errors.New("image url is invalid")
	}

	userID := uuid.New()
	if err := s.repo.CreateUser(ctx, repo.CreateUserInput{
		ID:        userID,
		ImageURL:  in.ImageURL,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
	}); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:        u.ID,
		ImageURL:  u.ImageURL,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
	}, nil
}

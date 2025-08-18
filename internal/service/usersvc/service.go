package usersvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
)

type userService interface {
	CreateUser(ctx context.Context, in user.User) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, in repo.User) error
	GetUser(ctx context.Context, id uuid.UUID) (repo.User, error)
}

type service struct {
	repo userRepository
}

func NewService(repo userRepository) *service {
	return &service{repo}
}

var _ userService = (*service)(nil)

func (s *service) CreateUser(ctx context.Context, in user.User) (uuid.UUID, error) {
	if in.FirstName == "" {
		return uuid.Nil, errors.New("first name is required")
	}
	if in.Username == "" {
		return uuid.Nil, errors.New("username is required")
	}

	userID := uuid.New()
	if err := s.repo.CreateUser(ctx, repo.User{
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
	if id == uuid.Nil {
		return user.User{}, errors.New("id is empty")
	}

	usr, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:        usr.ID,
		ImageURL:  usr.ImageURL,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Username:  usr.Username,
	}, nil
}

package usersvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
)

type userService interface {
	CreateUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, id uuid.UUID) (repo.User, error)
}

type service struct {
	repo userRepository
}

func NewService(repo userRepository) *service {
	return &service{repo}
}

var _ userService = (*service)(nil)

func (s *service) CreateUser(ctx context.Context, username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}

	if err := s.repo.CreateUser(ctx, username); err != nil {
		return err
	}

	return nil
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

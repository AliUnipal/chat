package inmemuserrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
)

func New(users ...repo.CreateUserInput) *repository {
	v := make(map[uuid.UUID]repo.CreateUserInput)
	for _, user := range users {
		v[user.ID] = user
	}

	return &repository{v}
}

type repository struct {
	users map[uuid.UUID]repo.CreateUserInput
}

func (r *repository) CreateUser(_ context.Context, in repo.CreateUserInput) error {
	if _, ok := r.users[in.ID]; ok {
		return errors.New("user already exists")
	}
	if in.ID == uuid.Nil {
		return errors.New("user id is required")
	}
	if in.FirstName == "" {
		return errors.New("first name is required")
	}
	if in.Username == "" {
		return errors.New("username is required")
	}

	r.users[in.ID] = in
	return nil
}

func (r *repository) GetUser(_ context.Context, id uuid.UUID) (repo.CreateUserInput, error) {
	user, ok := r.users[id]
	if !ok {
		return repo.CreateUserInput{}, errors.New("user does not exist")
	}

	return user, nil
}

package inmemuserrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
)

func New(users ...repo.User) *repository {
	v := make(map[uuid.UUID]repo.User)
	for _, user := range users {
		v[user.ID] = user
	}

	return &repository{v}
}

type repository struct {
	users map[uuid.UUID]repo.User
}

func (r *repository) CreateUser(_ context.Context, in repo.User) error {
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

func (r *repository) GetUser(_ context.Context, id uuid.UUID) (repo.User, error) {
	user, ok := r.users[id]
	if !ok {
		return repo.User{}, errors.New("user does not exist")
	}

	return user, nil
}

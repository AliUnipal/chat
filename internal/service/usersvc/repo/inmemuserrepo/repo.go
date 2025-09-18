package inmemuserrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/AliUnipal/chat/pkg/snapper"
	"github.com/google/uuid"
	"log"
)

// NOTE: The context is required for the snapper load.
func New(ctx context.Context) *repository {
	var data users

	s := snapper.NewFileSnapper[users]("user_data.json")
	d, err := s.Load(ctx)

	if err != nil {
		log.Println(err)
		data = make(users)
	} else {
		data = d
	}

	return &repository{data, s}
}

type users = map[uuid.UUID]repo.User

type repository struct {
	users   users
	snapper *snapper.FileSnapper[users]
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

	r.users[in.ID] = repo.User(in)
	return nil
}

func (r *repository) GetUser(_ context.Context, id uuid.UUID) (repo.User, error) {
	user, ok := r.users[id]
	if !ok {
		return repo.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.snapper.Snap(ctx, r.users)
}

package inmemuserrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/google/uuid"
)

type snapper interface {
	Snap(ctx context.Context, data userrepos.Data) error
	Load(ctx context.Context) (userrepos.Data, error)
}

func New(s snapper) *repository {
	return &repository{snapper: s}
}

type repository struct {
	users    userrepos.Data
	snapper  snapper
	isLoaded bool
}

func (r *repository) CreateUser(ctx context.Context, in userrepos.CreateUserInput) error {
	if !r.isLoaded {
		if err := r.Load(ctx); err != nil {
			return err
		}
	}

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

	r.users[in.ID] = userrepos.User(in)
	return nil
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (userrepos.User, error) {
	if !r.isLoaded {
		if err := r.Load(ctx); err != nil {
			return userrepos.User{}, err
		}
	}

	user, ok := r.users[id]
	if !ok {
		return userrepos.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r *repository) Load(ctx context.Context) error {
	data, err := r.snapper.Load(ctx)
	if err != nil {
		return err
	}

	r.users = data
	return nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.snapper.Snap(ctx, r.users)
}

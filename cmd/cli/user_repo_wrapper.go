package main

import (
	"context"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/google/uuid"
)

type userRepo interface {
	GetUser(ctx context.Context, id uuid.UUID) (userrepos.User, error)
}

type chatRepoUserRepoWrapper struct {
	chatRepo userRepo
}

func (c chatRepoUserRepoWrapper) GetUser(ctx context.Context, id uuid.UUID) (chatrepos.User, error) {
	u, err := c.chatRepo.GetUser(ctx, id)
	if err != nil {
		return chatrepos.User{}, err
	}

	return chatrepos.User{
		ID:        u.ID,
		ImageURL:  u.ImageURL,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
	}, nil
}

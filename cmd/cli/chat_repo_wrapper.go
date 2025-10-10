package main

import (
	"context"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos"
	"github.com/google/uuid"
)

type chatRepo interface {
	GetChat(ctx context.Context, id uuid.UUID) (*chatrepos.Chat, error)
}

type msgRepoChatRepoWrapper struct {
	msgRepo chatRepo
}

func (m *msgRepoChatRepoWrapper) GetChat(ctx context.Context, id uuid.UUID) (*msgrepos.Chat, error) {
	c, err := m.msgRepo.GetChat(ctx, id)
	if err != nil {
		return nil, err
	}

	return &msgrepos.Chat{
		ID:          c.ID,
		CurrentUser: msgrepos.User(c.CurrentUser),
		OtherUser:   msgrepos.User(c.OtherUser),
	}, nil
}

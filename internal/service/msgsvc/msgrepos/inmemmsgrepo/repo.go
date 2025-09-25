package inmemmsgrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos"
	"github.com/google/uuid"
)

type snapper interface {
	Snap(ctx context.Context, data msgrepos.Data) error
	Load(ctx context.Context) (msgrepos.Data, error)
}

type chatRepository interface {
	GetChat(ctx context.Context, id uuid.UUID) (*msgrepos.Chat, error)
}

func New(s snapper, chatRepo chatRepository) *repository {
	return &repository{
		chatRepo: chatRepo,
		snapper:  s,
	}
}

type repository struct {
	messages msgrepos.Data
	chatRepo chatRepository
	snapper  snapper
	isLoaded bool
}

func (r *repository) CreateMessage(ctx context.Context, in msgrepos.CreateMessageInput) error {
	if !r.isLoaded {
		return r.Load(ctx)
	}

	c, err := r.chatRepo.GetChat(ctx, in.ChatID)
	if err != nil {
		return err
	}
	if in.SenderID != c.CurrentUser.ID && in.SenderID != c.OtherUser.ID {
		return errors.New("user does not belong to this chat")
	}

	r.messages[in.ChatID] = append(r.messages[in.ChatID], msgrepos.Message{
		ID:          in.ID,
		SenderID:    in.SenderID,
		ChatID:      in.ChatID,
		Content:     in.Content,
		ContentType: in.ContentType,
		Timestamp:   in.Timestamp,
	})

	return nil
}

func (r *repository) GetMessages(ctx context.Context, chatID uuid.UUID) ([]msgrepos.Message, error) {
	if !r.isLoaded {
		return nil, r.Load(ctx)
	}

	msgs, ok := r.messages[chatID]
	if !ok {
		return nil, errors.New("chat does not exist")
	}

	return msgs, nil
}

func (r *repository) Load(ctx context.Context) error {
	data, err := r.snapper.Load(ctx)
	if err != nil {
		return err
	}

	r.messages = data
	return nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.snapper.Snap(ctx, r.messages)
}

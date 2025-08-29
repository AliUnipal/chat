package inmemmessagerepo

import (
	"context"
	"errors"
	chatrepo "github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo"
	"github.com/google/uuid"
)

func New(chatRepo chatRepository, msgs map[uuid.UUID][]repo.Message) *repository {
	return &repository{
		messages: msgs,
		chatRepo: chatRepo,
	}
}

type chatRepository interface {
	GetChat(ctx context.Context, id uuid.UUID) (chatrepo.Chat, error)
}

type repository struct {
	messages map[uuid.UUID][]repo.Message
	chatRepo chatRepository
}

func (r *repository) CreateMessage(ctx context.Context, in repo.CreateMessageInput) error {
	c, err := r.chatRepo.GetChat(ctx, in.ChatID)
	if err != nil {
		return err
	}
	if c.CurrentUser.ID != in.SenderID || c.OtherUser.ID != in.SenderID {
		return errors.New("user does not belong to this chat")
	}

	r.messages[in.ChatID] = append(r.messages[in.ChatID], repo.Message{
		ID:          in.ID,
		SenderID:    in.SenderID,
		ChatID:      in.ChatID,
		Content:     in.Content,
		ContentType: in.ContentType,
		Timestamp:   in.Timestamp,
	})

	return nil
}

func (r *repository) GetMessages(_ context.Context, chatID uuid.UUID) ([]repo.Message, error) {
	msgs, ok := r.messages[chatID]
	if !ok {
		return nil, errors.New("chat does not exist")
	}

	return msgs, nil
}

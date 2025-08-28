package inmemmessagerepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo"
	"github.com/google/uuid"
)

func New(chatRepo chatRepository, msgs []repo.Message) *repository {
	v := make(map[uuid.UUID][]repo.Message)
	for _, m := range msgs {
		_, ok := v[m.ChatID]
		if !ok {
			v[m.ChatID] = make([]repo.Message, 0)
		}
		v[m.ChatID] = append(v[m.ChatID], m)
	}

	return &repository{
		messages: v,
		chatRepo: chatRepo,
	}
}

type chatRepository interface {
	// TODO: Ask if it's better to do it like this, or to check by GetChatByID.
	// Probably this way is better for performance wise instead of fetching the whole chat.
	// Also do I need to make the return type as (bool, error)? Or unnessary
	CheckChatExist(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	messages map[uuid.UUID][]repo.Message
	chatRepo chatRepository
}

func (r *repository) CreateMessage(ctx context.Context, in repo.CreateMessageInput) error {
	err := r.chatRepo.CheckChatExist(ctx, in.ChatID)
	if err != nil {
		return err
	}

	//if _, ok := r.messages[in.ChatID]; !ok {
	//	r.messages[in.ChatID] = make([]repo.Message, 0)
	//}

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

func (r *repository) GetMessage(_ context.Context, id, chatID uuid.UUID) (repo.Message, error) {
	msgs, ok := r.messages[chatID]
	if !ok {
		return repo.Message{}, errors.New("chat does not exist")
	}

	for _, m := range msgs {
		if m.ID == id {
			return m, nil
		}
	}

	return repo.Message{}, errors.New("message not found")
}

func (r *repository) GetMessages(_ context.Context, chatID uuid.UUID) ([]repo.Message, error) {
	msgs, ok := r.messages[chatID]
	if !ok {
		return nil, errors.New("chat does not exist")
	}

	return msgs, nil
}

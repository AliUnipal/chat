package msgsvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos"
	"github.com/google/uuid"
	"time"
)

type MessageInput struct {
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType message.ContentType
}

// TODO: Ask about how the middleware/authorization for creating and getting message. And if it change the structure of the methods
type messageService interface {
	CreateMessage(ctx context.Context, in MessageInput) (uuid.UUID, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageRepository interface {
	CreateMessage(ctx context.Context, in msgrepos.CreateMessageInput) error
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]msgrepos.Message, error)
}

type service struct {
	repo messageRepository
}

var _ (messageService) = (*service)(nil)

func NewService(repo messageRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateMessage(ctx context.Context, in MessageInput) (uuid.UUID, error) {
	if in.Content == nil || len(in.Content) == 0 {
		return uuid.Nil, errors.New("content is empty")
	}
	if in.ChatID == uuid.Nil {
		return uuid.Nil, errors.New("chatID is empty")
	}
	if in.SenderID == uuid.Nil {
		return uuid.Nil, errors.New("senderID is empty")
	}

	id := uuid.New()

	if err := s.repo.CreateMessage(ctx, msgrepos.CreateMessageInput{
		ID:          id,
		SenderID:    in.SenderID,
		ChatID:      in.ChatID,
		Content:     in.Content,
		ContentType: in.ContentType,
		Timestamp:   time.Now().UTC(),
	}); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *service) GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error) {
	msgs, err := s.repo.GetMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}
	r := make([]message.Message, len(msgs))
	for i, m := range msgs {
		r[i] = message.Message{
			ID:          m.ID,
			SenderID:    m.SenderID,
			ChatID:      m.ChatID,
			Content:     m.Content,
			ContentType: m.ContentType,
			Timestamp:   m.Timestamp,
		}
	}

	return r, nil
}

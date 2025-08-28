package msgsvc

import (
	"context"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo"
	"github.com/google/uuid"
	"time"
)

type MessageInput struct {
	SenderID    uuid.UUID
	ChatID      uuid.UUID
	Content     []byte
	ContentType message.ContentType
	Timestamp   time.Time
}

// TODO: Ask about how the middleware/authorization for creating and getting message. And if it change the structure of the methods
type messageService interface {
	CreateMessage(ctx context.Context, in MessageInput) (uuid.UUID, error)
	GetMessage(ctx context.Context, id, chatID uuid.UUID) (message.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageRepository interface {
	CreateMessage(ctx context.Context, in repo.CreateMessageInput) error
	GetMessage(ctx context.Context, id, chatID uuid.UUID) (repo.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]repo.Message, error)
}

type service struct {
	repo messageRepository
}

var _ (messageService) = (*service)(nil)

func NewService(repo messageRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateMessage(ctx context.Context, in MessageInput) (uuid.UUID, error) {
	id := uuid.New()

	if err := s.repo.CreateMessage(ctx, repo.CreateMessageInput{
		ID:          id,
		SenderID:    in.SenderID,
		ChatID:      in.ChatID,
		Content:     in.Content,
		ContentType: in.ContentType,
		Timestamp:   in.Timestamp,
	}); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *service) GetMessage(ctx context.Context, id, chatID uuid.UUID) (message.Message, error) {
	m, err := s.repo.GetMessage(ctx, id, chatID)
	if err != nil {
		return message.Message{}, err
	}

	return message.Message{
		ID:          m.ID,
		SenderID:    m.SenderID,
		ChatID:      m.ChatID,
		Content:     m.Content,
		ContentType: m.ContentType,
		Timestamp:   m.Timestamp,
	}, nil
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

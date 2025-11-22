package messagecontroller

import (
	"context"

	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/google/uuid"
)

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (uuid.UUID, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

type messageController struct {
	messageSvc messageService
}

func New(messageSvc messageService) *messageController {
	return &messageController{messageSvc}
}

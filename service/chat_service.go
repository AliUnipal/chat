package service

import (
	"context"
	"github.com/AliUnipal/chat/model"
	"github.com/google/uuid"
)

// TODO:
// 1. Delete and recreate and complete ChatService interface
// 2. a) Implement the methods in the service struct
// 2. b) Update the ChatRepository interface to provide the necessary methods
// 2. c) Generate mocks for the updated interfaces
// 3. Write unit tests for the service methods

type ChatService interface {
	CreateChat(ctx context.Context, participantID *uuid.UUID, chatType model.ChatType) (model.Chat, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	EditChat(ctx context.Context, chatID uuid.UUID, body model.Chat) (model.Chat, error)
	InviteParticipant(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (model.User, error)
	RemoveParticipant(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error)
	TransferAdminRole(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (model.Chat, error) // is TransferOwnership better?
	ClearMessagesHistory(ctx context.Context, chatID uuid.UUID) error

	SendTextMessage(ctx context.Context, chatID uuid.UUID, content string) (model.Message, error)
	SendFileMessage(ctx context.Context, chatID uuid.UUID, content []byte) (model.Message, error)
	SendImageMessage(ctx context.Context, chatID uuid.UUID, content []byte) (model.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]model.Message, error)
}

type ChatRepository interface {
}

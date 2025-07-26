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
	DeleteChat(ctx context.Context, chatID uuid.UUID) (bool, error)
	EditChatName(ctx context.Context, chatID uuid.UUID, name string) (bool, error)
	EditChatImage(ctx context.Context, chatID uuid.UUID, imageURL string) (bool, error)
	InviteChatParticipant(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (model.User, error)
	RemoveChatParticipant(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (uuid.UUID, error)
	TransferAdminRole(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (model.Chat, error) // is TransferOwnership better?
}

type ChatRepository interface {
}

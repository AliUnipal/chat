package chatsvc

import (
	"context"
	"github.com/AliUnipal/chat/model"
	"github.com/google/uuid"
)

// TODO:
// 1. Delete and recreate and complete ChatService interface - Done
// 2. a) Implement the methods in the service struct - WIP
// 2. b) Update the ChatRepository interface to provide the necessary methods
// 2. c) Generate mocks for the updated interfaces
// 3. Write unit tests for the service methods

type ChatService interface {
	CreateDirectChat(ctx context.Context, ownerId uuid.UUID, participantIds *[]uuid.UUID) (uuid.UUID, error)
	InviteChatParticipants(ctx context.Context, chatID uuid.UUID, participantIds []uuid.UUID) ([]model.User, error)
	RemoveChatParticipant(ctx context.Context, chatID uuid.UUID, participantId uuid.UUID) (uuid.UUID, error)
	TransferAdminRole(ctx context.Context, chatID uuid.UUID, participantID uuid.UUID) (bool, error)
}

type ChatRepository interface {
}

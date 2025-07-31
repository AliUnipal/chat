package chatsvc

import (
	"context"
	"fmt"
	"github.com/AliUnipal/chat/internal/models"
	"github.com/google/uuid"
)

// TODO:
// 1. Delete and recreate and complete ChatService interface - Done
// 2. a) Implement the methods in the service struct - WIP
// 2. b) Update the ChatRepository interface to provide the necessary methods
// 2. c) Generate mocks for the updated interfaces
// 3. Write unit tests for the service methods

type ChatService interface {
	CreateDirectChat(ctx context.Context, ownerId uuid.UUID, participantId uuid.UUID) (uuid.UUID, error)
	// NOTE: This function in the future will either return all chat types "Group" and "Direct"
	// or they have to be seperate
	GetDirectChats(ctx context.Context, userId uuid.UUID) ([]models.DirectChat, error)
}

type ChatRepository interface {
	CreateDirectChat(ctx context.Context, chat models.DirectChat) (uuid.UUID, error)
	GetChatsByUserId(ctx context.Context, userId uuid.UUID) ([]models.DirectChat, error)
	// Ask about this if it should be in a different repo (User Repo) or a service.
	GetUserById(ctx context.Context, userId uuid.UUID) (models.User, error)
}

type UserContext interface {
	GetCurrentUserID(context context.Context) uuid.UUID // Ask if this should return error too?
}

type service struct {
	repo ChatRepository
	uc   UserContext
}

func NewService(repo ChatRepository, uc UserContext) *service {
	return &service{repo, uc}
}

func (s *service) CreateDirectChat(ctx context.Context, ownerId uuid.UUID, participantId uuid.UUID) (uuid.UUID, error) {
	userId := s.uc.GetCurrentUserID(ctx)
	participant, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return uuid.Nil, err
	}

	chat := models.DirectChat{
		ID:          uuid.New(),
		Name:        fmt.Sprintf("New chat with %s %s", participant.FirstName, participant.LastName),
		Admin:       userId,
		ImageURL:    "",
		Participant: participant,
	}
	chatID, err := s.repo.CreateDirectChat(ctx, chat)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func (s *service) GetDirectChats(ctx context.Context, userId uuid.UUID) ([]models.DirectChat, error) {
	chats, err := s.repo.GetChatsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

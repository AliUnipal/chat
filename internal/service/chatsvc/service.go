package chatsvc

import (
	"context"
	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	chatRepo "github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	userRepo "github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
)

// NOTE: Something that gives the pointer of the chat to add messages <- Repo, Dependency of chat repo.
// Dependency for msgsvc
// type chatRepository interface {
//	 GetChatByUser(ctx context.Context, userID uuid.UUID, userTwo uuid) (*repo.Chat, error)
// }

// TODO:
// 1. Delete and recreate and complete ChatService interface - Done
// 2. a) Implement the methods in the service struct - WIP
// 2. b) Update the ChatRepository interface to provide the necessary methods
// 2. c) Generate mocks for the updated interfaces
// 3. Write unit tests for the service methods

type chatService interface {
	CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) error
	GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error)
}

type chatRepository interface {
	CreateChat(ctx context.Context, chat chatRepo.CreateChatInput) error
	GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*chatRepo.Chat, error)
}

type userRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (userRepo.User, error)
}

var _ chatService = (*service)(nil)

func NewService(chatRepo chatRepository, userRepo userRepository) *service {
	return &service{chatRepo, userRepo}
}

type service struct {
	chatRepo chatRepository
	userRepo userRepository
}

func (s *service) CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) error {
	u1, err := s.userRepo.GetUser(ctx, currentUserID)
	if err != nil {
		return err
	}
	u2, err := s.userRepo.GetUser(ctx, otherUserID)
	if err != nil {
		return err
	}

	if err := s.chatRepo.CreateChat(ctx, chatRepo.CreateChatInput{
		CurrentUser: chatRepo.User{
			ID:        u1.ID,
			ImageURL:  u1.ImageURL,
			FirstName: u1.FirstName,
			LastName:  u1.LastName,
			Username:  u1.Username,
		},
		OtherUser: chatRepo.User{
			ID:        u2.ID,
			ImageURL:  u2.ImageURL,
			FirstName: u2.FirstName,
			LastName:  u2.LastName,
			Username:  u2.Username,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error) {
	c, err := s.chatRepo.GetChatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	chats := make([]chat.Chat, len(c))
	for i, c := range c {
		// TODO: Move the messages into its own service.
		messages := make([]message.Message, len(c.Messages))
		for j, m := range c.Messages {
			messages[j] = message.Message{
				ID:          m.ID,
				SenderID:    m.SenderID,
				Content:     m.Content,
				ContentType: message.TextContentType,
				Timestamp:   m.Timestamp,
			}
		}

		chats[i] = chat.Chat{
			ID: c.ID,
			CurrentUser: user.User{
				ID:        c.CurrentUser.ID,
				ImageURL:  c.CurrentUser.ImageURL,
				FirstName: c.CurrentUser.FirstName,
				LastName:  c.CurrentUser.LastName,
				Username:  c.CurrentUser.Username,
			},
			OtherUser: user.User{
				ID:        c.OtherUser.ID,
				ImageURL:  c.OtherUser.ImageURL,
				FirstName: c.OtherUser.FirstName,
				LastName:  c.OtherUser.LastName,
				Username:  c.OtherUser.Username,
			},
			Messages: messages,
		}
	}

	return chats, nil
}

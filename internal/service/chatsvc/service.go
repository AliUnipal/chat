package chatsvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
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
	CreateChat(ctx context.Context, chat repo.CreateChatInput) error
	GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*repo.Chat, error)
}

var _ chatService = (*service)(nil)

func NewService(repo chatRepository) *service {
	return &service{repo}
}

type service struct {
	repo chatRepository
}

func (s *service) CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) error {
	if currentUserID == uuid.Nil {
		return errors.New("current user ID is missing")
	}
	if otherUserID == uuid.Nil {
		return errors.New("other user ID is missing")
	}

	if err := s.repo.CreateChat(ctx, repo.CreateChatInput{
		UserOneID: currentUserID,
		UserTwoID: otherUserID,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error) {
	c, err := s.repo.GetChatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	chats := make([]chat.Chat, len(c))
	for i, c := range c {
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

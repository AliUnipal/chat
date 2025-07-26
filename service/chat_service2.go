package service

import (
	"context"
	"github.com/AliUnipal/chat/model"
	"github.com/google/uuid"
)

// import (
//
//	"context"
//	"fmt"
//	"slices"
//
//	"github.com/AliUnipal/chat/model"
//	"github.com/google/uuid"
//
// )
//
// // TODO:
// // 1. Delete and recreate and complete ChatService interface
// // 2. a) Implement the methods in the service struct
// // 2. b) Update the ChatRepository interface to provide the necessary methods
// // 2. c) Generate mocks for the updated interfaces
// // 3. Write unit tests for the service methods
type ChatService2 interface {
	CreateChat2(ctx context.Context, participantID uuid.UUID) (uuid.UUID, error)
	GetMessages2(ctx context.Context, chatID uuid.UUID) ([]model.Message, error)
	SendTextMessage2(ctx context.Context, chatID uuid.UUID, content string) (uuid.UUID, error)
}

type ChatWithMessages2 struct {
	Chat     model.Chat
	Messages []model.Message
}

type ChatRepository2 interface {
	CreateChat2(ctx context.Context, chat model.Chat) (uuid.UUID, error)
	GetChatByID2(ctx context.Context, chatID uuid.UUID) (ChatWithMessages2, error)
	SendMessage2(ctx context.Context, message model.Message) (uuid.UUID, error)
}

//
//type UserContext2 interface {
//	GetCurrentUserID(context.Context) uuid.UUID
//}
//
//func NewService2(repo ChatRepository, uc UserContext2) *service2 {
//	return &service2{repo, uc}
//}
//
//type service2 struct {
//	repo ChatRepository
//	uc   UserContext2
//}
//
//func (s *service2) GetMessages2(ctx context.Context, chatID uuid.UUID) ([]model.Message, error) {
//	chatWithMessages, err := s.repo.GetChatByID2(ctx, chatID)
//	if err != nil {
//		return nil, err
//	}
//
//	userID := s.uc.GetCurrentUserID(ctx)
//	if !slices.Contains(chatWithMessages.Chat2.Members, userID) && chatWithMessages.Chat.Admin.String() != userID.String() {
//		return nil, fmt.Errorf("user %s is not a member of chat %s", userID.String(), chatID.String())
//	}
//
//	return chatWithMessages.Messages2, nil
//}
//
//func (s *service2) CreateChat2(ctx context.Context, participantID uuid.UUID) (uuid.UUID, error) {
//	userID := s.uc.GetCurrentUserID(ctx)
//	chat := model.Chat{
//		ID:       uuid.New(),
//		Type:     model.DirectChatType,
//		Admin:    userID,
//		Name:     fmt.Sprintf("Chat between %s and %s", userID.String(), participantID.String()),
//		ImageURL: "",
//		Members:  []uuid.UUID{userID, participantID},
//	}
//
//	chatID, err := s.repo.CreateChat2(ctx, chat)
//	if err != nil {
//		return uuid.Nil, err
//	}
//
//	return chatID, nil
//}
//
//func (s *service2) SendTextMessage2(ctx context.Context, chatID uuid.UUID, content string) (uuid.UUID, error) {
//	senderID := s.uc.GetCurrentUserID(ctx)
//	message := model.Message{
//		ID:          uuid.New(),
//		SenderID:    senderID,
//		ChatID:      chatID,
//		Content:     []byte(content),
//		ContentType: model.TextContentType,
//	}
//
//	messageID, err := s.repo.SendMessage2(ctx, message)
//	if err != nil {
//		return uuid.Nil, err
//	}
//
//	return messageID, nil
//}

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AliUnipal/chat/model"
	"github.com/AliUnipal/chat/service"
	"github.com/AliUnipal/chat/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetMessages_NonParticipantMustFail(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	currentUserID := uuid.New()
	senderID := uuid.New()
	participantID := uuid.New()
	expectedMessages := []model.Message{
		{ID: uuid.New(), SenderID: senderID, ChatID: chatID, Content: []byte("Hello"), ContentType: model.TextContentType},
	}

	mockUserContext := mocks.NewUserContext(t)
	mockUserContext.EXPECT().GetCurrentUserID(mock.Anything).Return(currentUserID)

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatByID(mock.Anything, chatID).Return(service.ChatWithMessages{
		Chat: model.Chat{
			ID:       chatID,
			Type:     model.DirectChatType,
			Admin:    uuid.New(),
			Name:     "Something",
			ImageURL: "",
			Members:  []uuid.UUID{senderID, participantID},
		},
		Messages: expectedMessages,
	}, nil)

	service := service.NewService(mockRepo, mockUserContext)

	messages, err := service.GetMessages(ctx, chatID)
	if err == nil {
		t.Log("expected error, got nil")
		t.Fail()
	}

	if len(messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(messages))
	}
}

func TestGetMessages_AdminGetsMessages(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	currentUserID := uuid.New()
	senderID := uuid.New()
	participantID := uuid.New()
	expectedMessages := []model.Message{
		{ID: uuid.New(), SenderID: senderID, ChatID: chatID, Content: []byte("Hello"), ContentType: model.TextContentType},
	}

	mockUserContext := mocks.NewUserContext(t)
	mockUserContext.EXPECT().GetCurrentUserID(mock.Anything).Return(currentUserID)

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatByID(mock.Anything, chatID).Return(service.ChatWithMessages{
		Chat: model.Chat{
			ID:       chatID,
			Type:     model.DirectChatType,
			Admin:    currentUserID,
			Name:     "Something",
			ImageURL: "",
			Members:  []uuid.UUID{senderID, participantID},
		},
		Messages: expectedMessages,
	}, nil)

	service := service.NewService(mockRepo, mockUserContext)

	messages, err := service.GetMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(messages) != len(expectedMessages) {
								t.Fatalf("expected %d messages, got %d", len(expectedMessages), len(messages))
	}

	for i, msg := range messages {
		if msg.ID != expectedMessages[i].ID || msg.SenderID != expectedMessages[i].SenderID || msg.ChatID != expectedMessages[i].ChatID || string(msg.Content) != string(expectedMessages[i].Content) || msg.ContentType != expectedMessages[i].ContentType {
			t.Errorf("expected message %v, got %v", expectedMessages[i], msg)
		}
	}
}

func TestGetMessages_GetsMessages(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	currentUserID := uuid.New()
	participantID := uuid.New()
	expectedMessages := []model.Message{
		{ID: uuid.New(), SenderID: currentUserID, ChatID: chatID, Content: []byte("Hello"), ContentType: model.TextContentType},
		{ID: uuid.New(), SenderID: participantID, ChatID: chatID, Content: []byte("Hi"), ContentType: model.TextContentType},
	}

	mockUserContext := mocks.NewUserContext(t)
	mockUserContext.EXPECT().GetCurrentUserID(mock.Anything).Return(currentUserID)

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatByID(mock.Anything, chatID).Return(service.ChatWithMessages{
		Chat: model.Chat{
			ID:       chatID,
			Type:     model.DirectChatType,
			Admin:    currentUserID,
			Name:     "Something",
			ImageURL: "",
			Members:  []uuid.UUID{currentUserID, participantID},
		},
		Messages: expectedMessages,
	}, nil)

	service := service.NewService(mockRepo, mockUserContext)

	messages, err := service.GetMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(messages) != len(expectedMessages) {
		t.Fatalf("expected %d messages, got %d", len(expectedMessages), len(messages))
	}

	for i, msg := range messages {
		if msg.ID != expectedMessages[i].ID || msg.SenderID != expectedMessages[i].SenderID || msg.ChatID != expectedMessages[i].ChatID || string(msg.Content) != string(expectedMessages[i].Content) || msg.ContentType != expectedMessages[i].ContentType {
			t.Errorf("expected message %v, got %v", expectedMessages[i], msg)
		}
	}
}

func TestGetMessages_ReturnsErrors(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatByID(mock.Anything, chatID).Return(service.ChatWithMessages{}, errors.New("chat not found"))

	mockUserContext := mocks.NewUserContext(t)
	service := service.NewService(mockRepo, mockUserContext)

	if _, err := service.GetMessages(ctx, chatID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

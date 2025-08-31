package msgsvc_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/AliUnipal/chat/internal/service/msgsvc/mocks"
	"github.com/AliUnipal/chat/internal/service/msgsvc/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateMessage_ReturnID(t *testing.T) {
	ctx := context.Background()
	input := msgsvc.MessageInput{
		SenderID:    uuid.New(),
		ChatID:      uuid.New(),
		Content:     []byte("Hello Hello"),
		ContentType: message.TextContentType,
	}

	mockRepo := mocks.NewMessageRepository(t)
	mockRepo.EXPECT().CreateMessage(mock.Anything, mock.MatchedBy(func(r repo.CreateMessageInput) bool {
		return r.ID != uuid.Nil &&
			r.SenderID == input.SenderID &&
			r.ChatID == input.ChatID &&
			bytes.Equal(r.Content, input.Content) &&
			r.ContentType == input.ContentType
	})).Return(nil)

	service := msgsvc.NewService(mockRepo)

	id, err := service.CreateMessage(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("expected id got %v", id)
	}
}

func TestCreateMessage_ReturnErrorOnEmptyContent(t *testing.T) {
	ctx := context.Background()
	input := msgsvc.MessageInput{
		SenderID:    uuid.New(),
		ChatID:      uuid.New(),
		Content:     []byte(""),
		ContentType: message.TextContentType,
	}

	mockRepo := mocks.NewMessageRepository(t)

	service := msgsvc.NewService(mockRepo)
	if _, err := service.CreateMessage(ctx, input); err == nil {
		t.Fatalf("expected error got %v", err)
	}
}

func TestCreateMessage_ReturnErrorOnNilSenderID(t *testing.T) {
	ctx := context.Background()
	input := msgsvc.MessageInput{
		SenderID:    uuid.Nil,
		ChatID:      uuid.New(),
		Content:     []byte("Hello Hello"),
		ContentType: message.TextContentType,
	}

	mockRepo := mocks.NewMessageRepository(t)

	service := msgsvc.NewService(mockRepo)
	if _, err := service.CreateMessage(ctx, input); err == nil {
		t.Fatalf("expected error got %v", err)
	}
}

func TestCreateMessage_ReturnErrorOnNilChatID(t *testing.T) {
	ctx := context.Background()
	input := msgsvc.MessageInput{
		SenderID:    uuid.New(),
		ChatID:      uuid.Nil,
		Content:     []byte("Hello Hello"),
		ContentType: message.TextContentType,
	}

	mockRepo := mocks.NewMessageRepository(t)

	service := msgsvc.NewService(mockRepo)
	if _, err := service.CreateMessage(ctx, input); err == nil {
		t.Fatalf("expected error got %v", err)
	}
}

func TestCreateMessage_ReturnError(t *testing.T) {
	ctx := context.Background()
	input := msgsvc.MessageInput{
		SenderID:    uuid.New(),
		ChatID:      uuid.New(),
		Content:     []byte("Hello Hello"),
		ContentType: message.TextContentType,
	}

	mockRepo := mocks.NewMessageRepository(t)
	mockRepo.EXPECT().CreateMessage(mock.Anything, mock.MatchedBy(func(r repo.CreateMessageInput) bool {
		return r.ID != uuid.Nil &&
			r.SenderID == input.SenderID &&
			r.ChatID == input.ChatID &&
			bytes.Equal(r.Content, input.Content) &&
			r.ContentType == input.ContentType
	})).Return(errors.New("error"))

	service := msgsvc.NewService(mockRepo)
	if _, err := service.CreateMessage(ctx, input); err == nil {
		t.Fatalf("expected error got %v", err)
	}
}

func TestGetMessages_ReturnMessages(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	expectedMessages := []message.Message{
		{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      uuid.New(),
			Content:     []byte("Hello 1"),
			ContentType: message.TextContentType,
			Timestamp:   time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      uuid.New(),
			Content:     []byte("Hello 2"),
			ContentType: message.TextContentType,
			Timestamp:   time.Date(2009, time.November, 15, 23, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      uuid.New(),
			Content:     []byte("Hello 3"),
			ContentType: message.TextContentType,
			Timestamp:   time.Date(2009, time.November, 17, 23, 0, 0, 0, time.UTC),
		},
	}
	repoExpectedMessage := make([]repo.Message, len(expectedMessages))
	for i, msg := range expectedMessages {
		repoExpectedMessage[i] = repo.Message{
			ID:          msg.ID,
			SenderID:    msg.SenderID,
			ChatID:      msg.ChatID,
			Content:     msg.Content,
			ContentType: msg.ContentType,
			Timestamp:   msg.Timestamp,
		}
	}

	mockRepo := mocks.NewMessageRepository(t)
	mockRepo.EXPECT().GetMessages(mock.Anything, chatID).Return(repoExpectedMessage, nil)

	service := msgsvc.NewService(mockRepo)
	msgs, err := service.GetMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if len(msgs) != len(expectedMessages) {
		t.Fatalf("expected %d messages, got %d", len(expectedMessages), len(msgs))
	}

	for i, msg := range msgs {
		exM := expectedMessages[i]
		if msg.ID != exM.ID ||
			msg.SenderID != exM.SenderID ||
			msg.ChatID != exM.ChatID ||
			!bytes.Equal(msg.Content, exM.Content) ||
			msg.ContentType != exM.ContentType ||
			msg.Timestamp != exM.Timestamp {
			t.Fatalf("expected message %v, got %v", exM, msg)
		}
	}
}

func TestGetMessages_ReturnError(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()

	mockRepo := mocks.NewMessageRepository(t)
	mockRepo.EXPECT().GetMessages(mock.Anything, chatID).Return(nil, errors.New("error"))

	service := msgsvc.NewService(mockRepo)
	if _, err := service.GetMessages(ctx, chatID); err == nil {
		t.Fatalf("expected error got %v", err)
	}
}

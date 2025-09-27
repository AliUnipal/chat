package inmemmsgrepo_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos/inmemmsgrepo"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos/inmemmsgrepo/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func Test_CreateMessageReturnSuccess(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	senderID := uuid.New()
	in := msgrepos.CreateMessageInput{
		ID:          uuid.New(),
		SenderID:    senderID,
		ChatID:      chatID,
		Content:     nil,
		ContentType: 0,
		Timestamp:   time.Now(),
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{}, nil)
	mockChatRepo := mocks.NewChatRepository(t)
	mockChatRepo.EXPECT().GetChat(mock.Anything, chatID).Return(&msgrepos.Chat{
		ID: chatID,
		CurrentUser: msgrepos.User{
			ID: senderID,
		},
		OtherUser: msgrepos.User{},
	}, nil)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.CreateMessage(ctx, in); err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

// Don't judge the too long name.
func Test_CreateMessageReturnErrorOnUserNotBelongToChat(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	senderID := uuid.New()
	in := msgrepos.CreateMessageInput{
		ID:          uuid.New(),
		SenderID:    senderID,
		ChatID:      chatID,
		Content:     nil,
		ContentType: 0,
		Timestamp:   time.Now(),
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{}, nil)
	mockChatRepo := mocks.NewChatRepository(t)
	mockChatRepo.EXPECT().GetChat(mock.Anything, chatID).Return(&msgrepos.Chat{
		ID:          chatID,
		CurrentUser: msgrepos.User{},
		OtherUser:   msgrepos.User{},
	}, nil)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.CreateMessage(ctx, in); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func Test_CreateMessageReturnErrorOnSnapper(t *testing.T) {
	ctx := context.Background()
	in := msgrepos.CreateMessageInput{
		ID:          uuid.New(),
		SenderID:    uuid.New(),
		ChatID:      uuid.New(),
		Content:     nil,
		ContentType: 0,
		Timestamp:   time.Now(),
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{}, errors.New("error"))
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.CreateMessage(ctx, in); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func Test_GetMessagesReturnMsgs(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()
	expectedMsgs := []msgrepos.Message{
		msgrepos.Message{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      chatID,
			Content:     []byte("test 1"),
			ContentType: 0,
			Timestamp:   time.Now().Add(time.Minute * 100),
		},
		msgrepos.Message{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      chatID,
			Content:     []byte("test 2"),
			ContentType: 0,
			Timestamp:   time.Now().Add(time.Minute * 200),
		},
		msgrepos.Message{
			ID:          uuid.New(),
			SenderID:    uuid.New(),
			ChatID:      chatID,
			Content:     []byte("test 3"),
			ContentType: 0,
			Timestamp:   time.Now().Add(time.Minute * 300),
		},
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{
		chatID: expectedMsgs,
	}, nil)
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)
	msgs, err := repo.GetMessages(ctx, chatID)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if len(msgs) != len(expectedMsgs) {
		t.Fatalf("expected %d messages but got %d", len(expectedMsgs), len(msgs))
	}

	for i, m := range msgs {
		msg := expectedMsgs[i]

		if !reflect.DeepEqual(m, msg) {
			t.Fatalf("expected %v but got %v", msg, m)
		}
	}
}

func Test_GetMessagesReturnErrorChatNotExist(t *testing.T) {
	ctx := context.Background()
	chatID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{}, nil)
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)
	if _, err := repo.GetMessages(ctx, chatID); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_Load(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID][]msgrepos.Message{}, nil)
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.Load(ctx); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func Test_LoadReturnError(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, errors.New("error"))
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.Load(ctx); err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func Test_Close(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Snap(mock.Anything, mock.Anything).Return(nil)
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.Close(ctx); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func Test_CloseReturnError(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Snap(mock.Anything, mock.Anything).Return(errors.New("error"))
	mockChatRepo := mocks.NewChatRepository(t)

	repo := inmemmsgrepo.New(mockSnapper, mockChatRepo)

	if err := repo.Close(ctx); err == nil {
		t.Fatal("expected an error but got nil")
	}
}

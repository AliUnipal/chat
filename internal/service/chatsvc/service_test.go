package chatsvc_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc"
	"github.com/AliUnipal/chat/internal/service/chatsvc/mocks"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

// AAA - Arrange Act Assert

func TestGetChats_ReturnChats(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserOneID := uuid.New()
	expectedChats := []*repo.Chat{
		&repo.Chat{
			ID: userID,
			CurrentUser: repo.User{
				ID:        userID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			OtherUser: repo.User{
				ID:        otherUserOneID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			Messages: nil,
		},
		&repo.Chat{
			ID: userID,
			CurrentUser: repo.User{
				ID:        userID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			OtherUser: repo.User{
				ID:        uuid.UUID{},
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			Messages: nil,
		},
	}

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(expectedChats, nil)

	service := chatsvc.NewService(mockRepo)
	chats, err := service.GetChats(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(chats) != len(expectedChats) {
		t.Fatalf("expected %d chats, got %d", len(expectedChats), len(chats))
	}

	for i, c := range chats {
		if c.ID != expectedChats[i].ID || c.CurrentUser.ID != expectedChats[i].CurrentUser.ID || c.OtherUser.ID != expectedChats[i].OtherUser.ID {
			t.Errorf("expected chat \n%v, got \n%v", expectedChats[i], c)
		}
	}
}

func TestGetChats_ReturnErrors(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(nil, errors.New("not found"))

	service := chatsvc.NewService(mockRepo)
	if _, err := service.GetChats(ctx, userID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateChat_ReturnsErrors(t *testing.T) {
	ctx := context.Background()
	currentUserID := uuid.New()
	otherUserID := uuid.New()
	chatInput := repo.CreateChatInput{
		ID:        uuid.UUID{},
		UserOneID: currentUserID,
		UserTwoID: otherUserID,
	}

	mockRepo := mocks.NewChatRepository(t)
	mockRepo.EXPECT().CreateChat(mock.Anything, chatInput).Return(errors.New("error"))

	service := chatsvc.NewService(mockRepo)

	if err := service.CreateChat(ctx, currentUserID, otherUserID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

package chatsvc_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc"
	"github.com/AliUnipal/chat/internal/service/chatsvc/mocks"
	chatRepo "github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	"github.com/google/uuid"
	"testing"
)

// AAA - Arrange Act Assert

func TestGetChats_ReturnChats(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserOneID := uuid.New()
	expectedChats := []*chatRepo.Chat{
		&chatRepo.Chat{
			ID: userID,
			CurrentUser: chatRepo.User{
				ID:        userID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			OtherUser: chatRepo.User{
				ID:        otherUserOneID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			Messages: nil,
		},
		&chatRepo.Chat{
			ID: userID,
			CurrentUser: chatRepo.User{
				ID:        userID,
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			OtherUser: chatRepo.User{
				ID:        uuid.UUID{},
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
			Messages: nil,
		},
	}

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(expectedChats, nil)
	userMockRepo := mocks.NewUserRepository(t)

	service := chatsvc.NewService(chatMockRepo, userMockRepo)
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

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(nil, errors.New("not found"))
	userMockRepo := mocks.NewUserRepository(t)

	service := chatsvc.NewService(chatMockRepo, userMockRepo)
	if _, err := service.GetChats(ctx, userID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TODO: Fix unit tests.

//func TestCreateChat_ReturnsErrorsOnEmptyUserOneID(t *testing.T) {
//	ctx := context.Background()
//	currentUserID := uuid.New()
//	otherUserID := uuid.New()
//	chatInput := chatRepo.CreateChatInput{
//		ID: uuid.UUID{},
//		CurrentUser: chatRepo.User{
//			ID:        currentUserID,
//			ImageURL:  "",
//			FirstName: "Test 1",
//			LastName:  "",
//			Username:  "123",
//		},
//		OtherUser: chatRepo.User{
//			ID:        otherUserID,
//			ImageURL:  "",
//			FirstName: "Test 2",
//			LastName:  "",
//			Username:  "987",
//		},
//	}
//
//	chatMockRepo := mocks.NewChatRepository(t)
//	chatMockRepo.EXPECT().CreateChat(mock.Anything, chatInput).Return(errors.New("error"))
//	userMockRepo := mocks.NewUserRepository(t)
//	userMockRepo.EXPECT().GetUser(mock.Anything, currentUserID).Return(userRepo.User{}, errors.New("error"))
//
//	service := chatsvc.NewService(chatMockRepo, userMockRepo)
//
//	if err := service.CreateChat(ctx, currentUserID, otherUserID); err == nil {
//		t.Fatal("expected error, got nil")
//	}
//}

//func TestCreateChat_ReturnsErrors(t *testing.T) {
//	ctx := context.Background()
//	currentUserID := uuid.New()
//	otherUserID := uuid.New()
//	chatInput := chatRepo.CreateChatInput{
//		ID: uuid.UUID{},
//		CurrentUser: chatRepo.User{
//			ID:        currentUserID,
//			ImageURL:  "",
//			FirstName: "Test 1",
//			LastName:  "",
//			Username:  "123",
//		},
//		OtherUser: chatRepo.User{
//			ID:        otherUserID,
//			ImageURL:  "",
//			FirstName: "Test 2",
//			LastName:  "",
//			Username:  "987",
//		},
//	}
//
//	chatMockRepo := mocks.NewChatRepository(t)
//	chatMockRepo.EXPECT().CreateChat(mock.Anything, chatInput).Return(errors.New("error"))
//	userMockRepo := mocks.NewUserRepository(t)
//	userMockRepo.EXPECT().GetUser(mock.Anything, currentUserID).Return(userRepo.User{}, nil)
//	userMockRepo.EXPECT().GetUser(mock.Anything, otherUserID).Return(userRepo.User{}, nil)
//
//	service := chatsvc.NewService(chatMockRepo, userMockRepo)
//
//	if err := service.CreateChat(ctx, currentUserID, otherUserID); err == nil {
//		t.Fatal("expected error, got nil")
//	}
//}

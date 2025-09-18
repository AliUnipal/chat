package chatsvc_test

import (
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
	ctx := t.Context()
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
		},
	}

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(expectedChats, nil)

	service := chatsvc.NewService(chatMockRepo)
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

func TestGetChats_ReturnError(t *testing.T) {
	ctx := t.Context()
	userID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().GetChatsByUser(ctx, userID).Return(nil, errors.New("not found"))

	service := chatsvc.NewService(chatMockRepo)
	if _, err := service.GetChats(ctx, userID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateChat_ReturnID(t *testing.T) {
	ctx := t.Context()
	currentUserID := uuid.New()
	otherUserID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().CreateChat(mock.Anything, mock.MatchedBy(func(c repo.CreateChatInput) bool {
		return c.ID != uuid.Nil &&
			c.CurrentUserID == currentUserID &&
			c.OtherUserID == otherUserID
	})).Return(nil)

	service := chatsvc.NewService(chatMockRepo)

	id, err := service.CreateChat(ctx, currentUserID, otherUserID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == uuid.Nil {
		t.Fatal("expected id to be returned")
	}
}

func TestCreateChat_ReturnErrorOnEmptyUserOne(t *testing.T) {
	ctx := t.Context()
	otherUserID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().CreateChat(mock.Anything, mock.MatchedBy(func(c repo.CreateChatInput) bool {
		return c.ID != uuid.Nil &&
			c.CurrentUserID == uuid.Nil &&
			c.OtherUserID == otherUserID
	})).Return(errors.New("User one ID missing."))

	service := chatsvc.NewService(chatMockRepo)

	if _, err := service.CreateChat(ctx, uuid.Nil, otherUserID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateChat_ReturnErrorOnEmptyUserTwo(t *testing.T) {
	ctx := t.Context()
	currentUserID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().CreateChat(mock.Anything, mock.MatchedBy(func(c repo.CreateChatInput) bool {
		return c.ID != uuid.Nil &&
			c.CurrentUserID == currentUserID &&
			c.OtherUserID == uuid.Nil
	})).Return(errors.New("User one ID missing."))

	service := chatsvc.NewService(chatMockRepo)

	if _, err := service.CreateChat(ctx, currentUserID, uuid.Nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateChat_ReturnErrorOnIdenticalIDs(t *testing.T) {
	ctx := t.Context()
	currentUserID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().CreateChat(mock.Anything, mock.MatchedBy(func(c repo.CreateChatInput) bool {
		return c.ID != uuid.Nil &&
			c.CurrentUserID == currentUserID &&
			c.OtherUserID == currentUserID
	})).Return(errors.New("error"))

	service := chatsvc.NewService(chatMockRepo)

	if _, err := service.CreateChat(ctx, currentUserID, currentUserID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateChat_ReturnError(t *testing.T) {
	ctx := t.Context()
	currentUserID := uuid.New()
	otherUserID := uuid.New()

	chatMockRepo := mocks.NewChatRepository(t)
	chatMockRepo.EXPECT().CreateChat(mock.Anything, mock.MatchedBy(func(c repo.CreateChatInput) bool {
		return c.ID != uuid.Nil &&
			c.CurrentUserID == currentUserID &&
			c.OtherUserID == otherUserID
	})).Return(errors.New("error"))

	service := chatsvc.NewService(chatMockRepo)

	if _, err := service.CreateChat(ctx, currentUserID, otherUserID); err == nil {
		t.Fatal("expected error, got nil")
	}
}

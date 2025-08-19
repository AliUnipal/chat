package usersvc

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc/mocks"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser_ReturnID(t *testing.T) {
	ctx := context.Background()
	userInput := user.User{
		ImageURL:  "",
		FirstName: "First User",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(u repo.User) bool {
		return u.FirstName == userInput.FirstName &&
			u.LastName == userInput.LastName &&
			u.Username == userInput.Username &&
			u.ImageURL == userInput.ImageURL
	})).Return(nil)

	service := NewService(mockRepo)

	id, err := service.CreateUser(ctx, userInput)
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}

	if id == uuid.Nil {
		t.Fatalf("Expected id %v got %v", userInput.ID, id)
	}
}

func TestCreateUser_ReturnErrorOnEmptyFirstName(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	userInput := user.User{
		ID:       userID,
		ImageURL: "",
		LastName: "Last Name",
		Username: "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	service := NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestCreateUser_ReturnErrorOnEmptyUsername(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	userInput := user.User{
		ID:        userID,
		ImageURL:  "",
		FirstName: "First Name",
		LastName:  "Last Name",
	}

	mockRepo := mocks.NewUserRepository(t)
	service := NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestCreateUser_ReturnError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	userInput := user.User{
		ID:        userID,
		ImageURL:  "",
		FirstName: "First Name",
		LastName:  "Last Name",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(errors.New("error"))

	service := NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestGetUser_ReturnUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	expectedUser := user.User{
		ID:        userID,
		ImageURL:  "",
		FirstName: "First Name",
		LastName:  "Last Name",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().GetUser(ctx, userID).Return(repo.User{
		ID:        expectedUser.ID,
		ImageURL:  expectedUser.ImageURL,
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Username:  expectedUser.Username,
	}, nil)
	service := NewService(mockRepo)

	usr, err := service.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}

	if usr != expectedUser {
		t.Fatalf("Expected user %v got %v", expectedUser, usr)
	}
}

func TestGetUser_ReturnErrorOnEmptyID(t *testing.T) {
	ctx := context.Background()

	mockRepo := mocks.NewUserRepository(t)
	service := NewService(mockRepo)

	if _, err := service.GetUser(ctx, uuid.Nil); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestGetUser_ReturnError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().GetUser(mock.Anything, userID).Return(repo.User{}, errors.New("error"))

	service := NewService(mockRepo)

	if _, err := service.GetUser(ctx, userID); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

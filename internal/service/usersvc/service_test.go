package usersvc_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/AliUnipal/chat/internal/service/usersvc/mocks"
	"github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUser_ReturnID(t *testing.T) {
	ctx := context.Background()
	userInput := usersvc.CreateUserInput{
		ImageURL:  "https://test.png",
		FirstName: "First CreateUserInput",
		LastName:  "Test",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(u repo.CreateUserInput) bool {
		return u.ID != uuid.Nil &&
			u.FirstName == userInput.FirstName &&
			u.LastName == userInput.LastName &&
			u.Username == userInput.Username &&
			u.ImageURL == userInput.ImageURL
	})).Return(nil)

	service := usersvc.NewService(mockRepo)

	id, err := service.CreateUser(ctx, userInput)
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("Expected user id got %v", id)
	}
}

func TestCreateUser_ReturnErrorOnEmptyFirstName(t *testing.T) {
	ctx := context.Background()
	userInput := usersvc.CreateUserInput{
		ImageURL: "https://test.png",
		LastName: "Last Name",
		Username: "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	service := usersvc.NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestCreateUser_ReturnErrorOnEmptyUsername(t *testing.T) {
	ctx := context.Background()
	userInput := usersvc.CreateUserInput{
		ImageURL:  "https://test.png",
		FirstName: "First Name",
		LastName:  "Last Name",
	}

	mockRepo := mocks.NewUserRepository(t)
	service := usersvc.NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestCreateUser_ReturnErrorOnInvalidImageURL(t *testing.T) {
	ctx := context.Background()
	userInput := usersvc.CreateUserInput{
		ImageURL:  "/test.png",
		FirstName: "First Name",
		LastName:  "Last Name",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	service := usersvc.NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestCreateUser_ReturnError(t *testing.T) {
	ctx := context.Background()
	userInput := usersvc.CreateUserInput{
		ImageURL:  "https://test.png",
		FirstName: "First Name",
		LastName:  "Last Name",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(errors.New("error"))

	service := usersvc.NewService(mockRepo)

	if _, err := service.CreateUser(ctx, userInput); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestGetUser_ReturnUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	expectedUser := user.User{
		ID:        userID,
		ImageURL:  "https://test.png",
		FirstName: "First Name",
		LastName:  "Last Name",
		Username:  "+97312345678",
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.EXPECT().GetUser(ctx, userID).Return(repo.CreateUserInput{
		ID:        expectedUser.ID,
		ImageURL:  expectedUser.ImageURL,
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Username:  expectedUser.Username,
	}, nil)
	service := usersvc.NewService(mockRepo)

	usr, err := service.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}

	if usr != expectedUser {
		t.Fatalf("Expected user %v got %v", expectedUser, usr)
	}
}

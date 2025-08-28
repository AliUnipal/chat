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
	userInput := CreateUserInput{
		ImageURL:  "https://s3....",
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

	service := NewService(mockRepo)

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
	userInput := CreateUserInput{
		ImageURL: "Image",
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
	userInput := CreateUserInput{
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
	userInput := CreateUserInput{
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
		ImageURL:  "https://test..",
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
	mockRepo.EXPECT().GetUser(mock.Anything, userID).Return(repo.CreateUserInput{}, errors.New("error"))

	service := NewService(mockRepo)

	if _, err := service.GetUser(ctx, userID); err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

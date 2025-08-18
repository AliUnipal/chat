package usersvc

import (
	"testing"
)

func TestCreateUser_ReturnID(t *testing.T) {
	//ctx := context.Background()
	//userInput := user.User{
	//	ImageURL:  "",
	//	FirstName: "First User",
	//	LastName:  "",
	//	Username:  "+97312345678",
	//}
	//
	//mockRepo := mocks.NewUserRepository(t)
	//mockRepo.EXPECT().CreateUser(mock.Anything, repo.User{
	//	ImageURL:  userInput.ImageURL,
	//	FirstName: userInput.FirstName,
	//	LastName:  userInput.LastName,
	//	Username:  userInput.Username,
	//}).Return(nil)
	//
	//service := NewService(mockRepo)
	//
	//_, err := service.CreateUser(ctx, userInput)
	//if err != nil {
	//	t.Fatalf("Expected no error got %v", err)
	//}
	//
	//if id != userInput.ID {
	//	t.Fatalf("Expected id %v got %v", userInput.ID, id)
	//}
}

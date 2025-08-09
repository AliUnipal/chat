package chatsvc_test

// AAA - Arrange Act Assert

//func TestCreateChat_ReturnId(t *testing.T) {
//	ctx := context.Background()
//
//}

//func TestCreateChat_ReturnsErrors(t *testing.T) {
//	ctx := context.Background()
//	currentUserId := uuid.New()
//	participantId := uuid.New()
//	chat := models.Chat{
//		ID:          uuid.UUID{},
//		Admin:       currentUserId,
//		Name:        "",
//		ImageURL:    "",
//		Participant: models.User{},
//	}
//
//	mockRepo := mocks.NewChatRepository(t)
//	mockRepo.EXPECT().CreateChat(mock.Anything, chat).Return(uuid.New(), errors.New("error"))
//
//	mockUserContext := mocks.NewUserContext(t)
//
//	service := chatsvc.NewService(mockRepo, mockUserContext)
//
//	if _, err := service.CreateChat(ctx, currentUserId, participantId); err != nil {
//		t.Fatal("expected error, got nil")
//	}
//}

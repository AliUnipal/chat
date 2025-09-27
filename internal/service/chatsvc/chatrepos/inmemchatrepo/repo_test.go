package inmemchatrepo_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos/inmemchatrepo"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos/inmemchatrepo/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func Test_CreateChatReturnSuccess(t *testing.T) {
	ctx := context.Background()
	cUsrID := uuid.New()
	oUsrID := uuid.New()
	in := chatrepos.CreateChatInput{
		ID:            uuid.New(),
		CurrentUserID: cUsrID,
		OtherUserID:   oUsrID,
	}
	expCurUsr := chatrepos.User{
		ID:        cUsrID,
		ImageURL:  "",
		FirstName: "Curr",
		LastName:  "Curr",
		Username:  "+9712345678",
	}
	expOUsr := chatrepos.User{
		ID:        oUsrID,
		ImageURL:  "",
		FirstName: "Other",
		LastName:  "Other",
		Username:  "+9712345679",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, cUsrID).Return(expCurUsr, nil)
	mockUserRepo.EXPECT().GetUser(mock.Anything, oUsrID).Return(expOUsr, nil)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if err := repo.CreateChat(ctx, in); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func Test_CreateChatReturnErrorOnSnapper(t *testing.T) {
	ctx := context.Background()
	in := chatrepos.CreateChatInput{}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, errors.New("Error"))
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)

	if err := repo.CreateChat(ctx, in); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_CreateChatReturnErrorOnChatExists(t *testing.T) {
	ctx := context.Background()
	cUsrID := uuid.New()
	oUsrID := uuid.New()
	in := chatrepos.CreateChatInput{
		ID:            uuid.New(),
		CurrentUserID: cUsrID,
		OtherUserID:   oUsrID,
	}
	ids := []string{cUsrID.String(), oUsrID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	expcChats := map[string]*chatrepos.Chat{
		id: &chatrepos.Chat{},
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{
		Chats:     expcChats,
		UserChats: make(map[uuid.UUID][]*chatrepos.Chat),
	}, nil)
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if err := repo.CreateChat(ctx, in); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_CreateChatReturnErrorOnNoCurUser(t *testing.T) {
	ctx := context.Background()
	cUsrID := uuid.New()
	oUsrID := uuid.New()
	in := chatrepos.CreateChatInput{
		ID:            uuid.New(),
		CurrentUserID: cUsrID,
		OtherUserID:   oUsrID,
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, cUsrID).Return(chatrepos.User{}, errors.New("doesn't exist"))

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if err := repo.CreateChat(ctx, in); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_CreateChatReturnErrorOnNoOtherUser(t *testing.T) {
	ctx := context.Background()
	cUsrID := uuid.New()
	oUsrID := uuid.New()
	in := chatrepos.CreateChatInput{
		ID:            uuid.New(),
		CurrentUserID: cUsrID,
		OtherUserID:   oUsrID,
	}
	expCurUsr := chatrepos.User{
		ID:        cUsrID,
		ImageURL:  "",
		FirstName: "Curr",
		LastName:  "Curr",
		Username:  "+9712345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, cUsrID).Return(expCurUsr, nil)
	mockUserRepo.EXPECT().GetUser(mock.Anything, oUsrID).Return(chatrepos.User{}, errors.New("doesn't exist"))

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if err := repo.CreateChat(ctx, in); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_GetChatReturnChat(t *testing.T) {
	ctx := context.Background()
	cUsrID := uuid.New()
	oUsrID := uuid.New()
	expChat := chatrepos.Chat{
		ID: uuid.New(),
		CurrentUser: chatrepos.User{
			ID: cUsrID,
		},
		OtherUser: chatrepos.User{
			ID: oUsrID,
		},
	}
	id := cUsrID.String() + "|" + oUsrID.String()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{
		Chats: map[string]*chatrepos.Chat{
			id: &expChat,
		},
		UserChats: make(map[uuid.UUID][]*chatrepos.Chat),
	}, nil)
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	c, err := repo.GetChat(ctx, expChat.ID)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if c.ID != expChat.ID || c.CurrentUser.ID != expChat.CurrentUser.ID || c.OtherUser.ID != expChat.OtherUser.ID {
		t.Fatalf("expected %v but got %v", expChat, c)
	}
}

func Test_GetChatReturnErrorNotFound(t *testing.T) {
	ctx := context.Background()
	cID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)

	if _, err := repo.GetChat(ctx, cID); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_GetChatReturnErrorOnSnapper(t *testing.T) {
	ctx := context.Background()
	cID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, errors.New("error"))
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)

	if _, err := repo.GetChat(ctx, cID); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func Test_GetChatsByUserReturnChats(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	expCurUsr := chatrepos.User{
		ID:        uID,
		ImageURL:  "",
		FirstName: "Curr",
		LastName:  "Curr",
		Username:  "+9712345678",
	}

	expectedChats := []*chatrepos.Chat{
		&chatrepos.Chat{
			ID:          uuid.New(),
			CurrentUser: expCurUsr,
			OtherUser: chatrepos.User{
				ID:        uuid.New(),
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
		},
		&chatrepos.Chat{
			ID:          uuid.New(),
			CurrentUser: expCurUsr,
			OtherUser: chatrepos.User{
				ID:        uuid.New(),
				ImageURL:  "",
				FirstName: "",
				LastName:  "",
				Username:  "",
			},
		},
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{
		Chats: map[string]*chatrepos.Chat{},
		UserChats: map[uuid.UUID][]*chatrepos.Chat{
			uID: expectedChats,
		},
	}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, uID).Return(expCurUsr, nil)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	cs, err := repo.GetChatsByUser(ctx, uID)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	for i, c := range cs {
		ec := expectedChats[i]
		if !reflect.DeepEqual(ec, c) {
			t.Fatalf("expected %+v but got %+v", ec, c)
		}
	}
}

func Test_GetChatsByUserReturnErrorNotFound(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	expCurUsr := chatrepos.User{
		ID:        uID,
		ImageURL:  "",
		FirstName: "Curr",
		LastName:  "Curr",
		Username:  "+9712345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, uID).Return(expCurUsr, nil)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if _, err := repo.GetChatsByUser(ctx, uID); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func Test_GetChatsByUserReturnErrorOnUserNotFound(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, nil)
	mockUserRepo := mocks.NewUserRepository(t)
	mockUserRepo.EXPECT().GetUser(mock.Anything, uID).Return(chatrepos.User{}, errors.New("error"))

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if _, err := repo.GetChatsByUser(ctx, uID); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func Test_GetChatsByUserReturnErrorOnSnapper(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(chatrepos.Data{}, errors.New("error"))
	mockUserRepo := mocks.NewUserRepository(t)

	repo := inmemchatrepo.New(mockSnapper, mockUserRepo)
	if _, err := repo.GetChatsByUser(ctx, uID); err == nil {
		t.Fatalf("expected error but got nil")
	}
}

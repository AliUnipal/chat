package inmemuserrepo_test

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos/inmemuserrepo"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos/inmemuserrepo/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_CreateUser(t *testing.T) {
	ctx := context.Background()
	in := userrepos.CreateUserInput{
		ID:        uuid.New(),
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID]userrepos.User{}, nil)

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

func Test_CreateUserReturnErrorSnapper(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	in := userrepos.CreateUserInput{
		ID:        userID,
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, errors.New("error"))

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_CreateUserReturnErrorAlreadyExists(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	in := userrepos.CreateUserInput{
		ID:        userID,
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID]userrepos.User{
		userID: userrepos.User{
			ID:        userID,
			ImageURL:  "",
			FirstName: "User 1",
			LastName:  "",
			Username:  "+97312345678",
		},
	}, nil)

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_CreateUserReturnErrorIDRequired(t *testing.T) {
	ctx := context.Background()
	in := userrepos.CreateUserInput{
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, nil)

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_CreateUserReturnErrorFirstNameRequired(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	in := userrepos.CreateUserInput{
		ID:       userID,
		ImageURL: "",
		LastName: "",
		Username: "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, nil)

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_CreateUserReturnErrorUsernameRequired(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	in := userrepos.CreateUserInput{
		ID:        userID,
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, nil)

	repo := inmemuserrepo.New(mockSnapper)

	if err := repo.CreateUser(ctx, in); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_GetUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	expUser := userrepos.User{
		ID:        userID,
		ImageURL:  "",
		FirstName: "User 1",
		LastName:  "",
		Username:  "+97312345678",
	}

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(map[uuid.UUID]userrepos.User{
		userID: expUser,
	}, nil)

	repo := inmemuserrepo.New(mockSnapper)
	user, err := repo.GetUser(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if user != expUser {
		t.Fatalf("expected user %v got %v", expUser, user)
	}
}

func Test_GetUserReturnErrorSnapper(t *testing.T) {
	ctx := context.Background()
	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, errors.New("error"))

	repo := inmemuserrepo.New(mockSnapper)
	if _, err := repo.GetUser(ctx, uuid.New()); err == nil {
		t.Fatal("expected error got nil")
	}
}

func Test_GetUserReturnUserDoesNotExist(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, nil)

	repo := inmemuserrepo.New(mockSnapper)
	if _, err := repo.GetUser(ctx, userID); err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestRepository_Load(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, nil)

	repo := inmemuserrepo.New(mockSnapper)
	if err := repo.Load(ctx); err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

func TestRepository_LoadReturnError(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Load(mock.Anything).Return(nil, errors.New("error"))

	repo := inmemuserrepo.New(mockSnapper)
	if err := repo.Load(ctx); err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestRepository_Close(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Snap(mock.Anything, mock.Anything).Return(nil)

	repo := inmemuserrepo.New(mockSnapper)
	if err := repo.Close(ctx); err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

func TestRepository_CloseReturnError(t *testing.T) {
	ctx := context.Background()

	mockSnapper := mocks.NewSnapper(t)
	mockSnapper.EXPECT().Snap(mock.Anything, mock.Anything).Return(errors.New("error"))

	repo := inmemuserrepo.New(mockSnapper)
	if err := repo.Close(ctx); err == nil {
		t.Fatal("expected error got nil")
	}
}

package inmemchatrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	userRepo "github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/google/uuid"
	"slices"
	"strings"
)

func New(userRepo userRepository) *repository {
	return &repository{
		make(map[string]*repo.Chat),
		make(map[uuid.UUID][]*repo.Chat),
		userRepo,
	}
}

type repository struct {
	chats     map[string]*repo.Chat
	userChats map[uuid.UUID][]*repo.Chat
	userRepo  userRepository
}

type userRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (userRepo.CreateUserInput, error)
}

func (r *repository) CreateChat(ctx context.Context, in repo.CreateChatInput) error {
	ids := []string{in.CurrentUserID.String(), in.OtherUserID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	if _, ok := r.chats[id]; ok {
		return errors.New("chat already exists")
	}

	cu, err := r.userRepo.GetUser(ctx, in.CurrentUserID)
	if err != nil {
		return err
	}
	ou, err := r.userRepo.GetUser(ctx, in.OtherUserID)
	if err != nil {
		return err
	}

	chat := &repo.Chat{
		ID:          in.ID,
		CurrentUser: repo.User(cu),
		OtherUser:   repo.User(ou),
	}

	r.chats[id] = chat
	r.userChats[in.CurrentUserID] = append(r.userChats[in.CurrentUserID], chat)
	r.userChats[in.OtherUserID] = append(r.userChats[in.OtherUserID], chat)

	return nil
}

func (r *repository) GetChatsByUser(_ context.Context, userID uuid.UUID) ([]*repo.Chat, error) {
	return r.userChats[userID], nil
}

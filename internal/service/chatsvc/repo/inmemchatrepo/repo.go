package inmemchatrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	"github.com/google/uuid"
	"slices"
	"strings"
)

func New() *repository {
	return &repository{
		make(map[string]*repo.Chat),
		make(map[uuid.UUID][]*repo.Chat),
	}
}

type repository struct {
	chats     map[string]*repo.Chat
	userChats map[uuid.UUID][]*repo.Chat
}

func (r *repository) CreateChat(_ context.Context, in repo.CreateChatInput) error {
	ids := []string{in.CurrentUser.ID.String(), in.OtherUser.ID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	if _, ok := r.chats[id]; ok {
		return errors.New("chat already exists")
	}

	chat := &repo.Chat{
		ID:          in.ID,
		CurrentUser: in.CurrentUser,
		OtherUser:   in.OtherUser,
	}

	r.chats[id] = chat
	r.userChats[in.CurrentUser.ID] = append(r.userChats[in.CurrentUser.ID], chat)
	r.userChats[in.OtherUser.ID] = append(r.userChats[in.OtherUser.ID], chat)

	return nil
}

func (r *repository) GetChatsByUser(_ context.Context, userID uuid.UUID) ([]*repo.Chat, error) {
	return r.userChats[userID], nil
}

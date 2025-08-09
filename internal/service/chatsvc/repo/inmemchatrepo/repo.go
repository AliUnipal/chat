package inmemchatrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	"github.com/google/uuid"
	"slices"
	"strings"
)

func New(users ...repo.User) *repository {
	m := make(map[uuid.UUID]repo.User)
	for _, user := range users {
		m[user.ID] = user
	}

	return &repository{
		make(map[string]*repo.Chat),
		make(map[uuid.UUID][]*repo.Chat),
		m,
	}
}

type repository struct {
	chats     map[string]*repo.Chat
	userChats map[uuid.UUID][]*repo.Chat
	users     map[uuid.UUID]repo.User
}

func (r *repository) CreateChat(_ context.Context, in repo.CreateChatInput) error {
	u1, ok := r.users[in.UserOneID]
	if !ok {
		return errors.New("user one does not exist")
	}
	u2, ok := r.users[in.UserTwoID]
	if !ok {
		return errors.New("user two does not exist")
	}
	ids := []string{in.UserOneID.String(), in.UserTwoID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	if _, ok := r.chats[id]; ok {
		return errors.New("chat already exists")
	}

	chat := &repo.Chat{
		ID:          in.ID,
		CurrentUser: u1,
		OtherUser:   u2,
	}

	r.chats[id] = chat
	r.userChats[u1.ID] = append(r.userChats[u1.ID], chat)
	r.userChats[u2.ID] = append(r.userChats[u2.ID], chat)

	return nil
}

func (r *repository) GetChatsByUser(_ context.Context, userID uuid.UUID) ([]*repo.Chat, error) {
	return r.userChats[userID], nil
}

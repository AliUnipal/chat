package inmemchatrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/google/uuid"
	"slices"
	"strings"
)

type snapper interface {
	Snap(ctx context.Context, data chatrepos.Data) error
	Load(ctx context.Context) (chatrepos.Data, error)
}

type userRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (chatrepos.User, error)
}

func New(s snapper, userRepo userRepository) *repository {
	return &repository{
		userRepo: userRepo,
		snapper:  s,
	}
}

type repository struct {
	data     chatrepos.Data
	userRepo userRepository
	snapper  snapper
	isLoaded bool
}

func (r *repository) CreateChat(ctx context.Context, in chatrepos.CreateChatInput) error {
	if !r.isLoaded {
		return r.Load(ctx)
	}

	ids := []string{in.CurrentUserID.String(), in.OtherUserID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	if _, ok := r.data.Chats[id]; ok {
		return errors.New("chat already exists")
	}
	if in.CurrentUserID == in.OtherUserID {
		return errors.New("user ids are matching, you cannot create a chat with yourself!")
	}

	cu, err := r.userRepo.GetUser(ctx, in.CurrentUserID)
	if err != nil {
		return err
	}
	ou, err := r.userRepo.GetUser(ctx, in.OtherUserID)
	if err != nil {
		return err
	}

	chat := &chatrepos.Chat{
		ID:          in.ID,
		CurrentUser: chatrepos.User(cu),
		OtherUser:   chatrepos.User(ou),
	}
	oppositeChat := &chatrepos.Chat{
		ID:          in.ID,
		CurrentUser: chatrepos.User(ou),
		OtherUser:   chatrepos.User(cu),
	}

	r.data.Chats[id] = chat
	r.data.UserChats[in.CurrentUserID] = append(r.data.UserChats[in.CurrentUserID], chat)
	r.data.UserChats[in.OtherUserID] = append(r.data.UserChats[in.OtherUserID], oppositeChat)

	return nil
}

func (r *repository) GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*chatrepos.Chat, error) {
	if !r.isLoaded {
		return nil, r.Load(ctx)
	}

	_, err := r.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.data.UserChats[userID], nil
}

func (r *repository) GetChat(ctx context.Context, id uuid.UUID) (*chatrepos.Chat, error) {
	if !r.isLoaded {
		return nil, r.Load(ctx)
	}

	for _, chat := range r.data.Chats {
		if chat.ID == id {
			return chat, nil
		}
	}

	return nil, errors.New("chat not found")
}

func (r *repository) Load(ctx context.Context) error {
	data, err := r.snapper.Load(ctx)
	if err != nil {
		return err
	}

	r.data.UserChats = data.UserChats
	r.data.Chats = data.Chats

	return nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.snapper.Snap(ctx, chatrepos.Data{
		Chats:     r.data.Chats,
		UserChats: r.data.UserChats,
	})
}

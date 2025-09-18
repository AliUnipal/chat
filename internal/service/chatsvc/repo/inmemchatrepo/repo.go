package inmemchatrepo

import (
	"context"
	"errors"
	"github.com/AliUnipal/chat/internal/service/chatsvc/repo"
	userRepo "github.com/AliUnipal/chat/internal/service/usersvc/repo"
	"github.com/AliUnipal/chat/pkg/snapper"
	"github.com/google/uuid"
	"log"
	"slices"
	"strings"
)

// NOTE: The context is required for the snapper load.
func New(ctx context.Context, userRepo userRepository) *repository {
	var chats data

	s := snapper.NewFileSnapper[data]("chats_data.json")
	d, err := s.Load(ctx)
	if err != nil {
		log.Println(err)
		chats = data{
			make(map[string]*repo.Chat),
			make(map[uuid.UUID][]*repo.Chat),
		}
	} else {
		chats = d
	}

	return &repository{
		chats,
		userRepo,
		s,
	}
}

type data struct {
	Chats     map[string]*repo.Chat      `json:"chats"`
	UserChats map[uuid.UUID][]*repo.Chat `json:"user_chats"`
}

type repository struct {
	data
	userRepo userRepository
	snapper  *snapper.FileSnapper[data]
}

type userRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (userRepo.User, error)
}

func (r *repository) CreateChat(ctx context.Context, in repo.CreateChatInput) error {
	ids := []string{in.CurrentUserID.String(), in.OtherUserID.String()}
	slices.Sort(ids)
	id := strings.Join(ids, "|")
	if _, ok := r.Chats[id]; ok {
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

	chat := &repo.Chat{
		ID:          in.ID,
		CurrentUser: repo.User(cu),
		OtherUser:   repo.User(ou),
	}
	oppositeChat := &repo.Chat{
		ID:          in.ID,
		CurrentUser: repo.User(ou),
		OtherUser:   repo.User(cu),
	}

	r.Chats[id] = chat
	r.UserChats[in.CurrentUserID] = append(r.UserChats[in.CurrentUserID], chat)
	r.UserChats[in.OtherUserID] = append(r.UserChats[in.OtherUserID], oppositeChat)

	return nil
}

func (r *repository) GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*repo.Chat, error) {
	_, err := r.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.UserChats[userID], nil
}

func (r *repository) GetChat(_ context.Context, id uuid.UUID) (*repo.Chat, error) {
	for _, chat := range r.Chats {
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

	r.UserChats = data.UserChats
	r.Chats = data.Chats

	return nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.snapper.Snap(ctx, data{
		Chats:     r.Chats,
		UserChats: r.UserChats,
	})
}

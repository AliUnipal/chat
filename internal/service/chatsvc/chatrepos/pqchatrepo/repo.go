package pqchatrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos/pqchatrepo/queries"
	"github.com/google/uuid"
)

//go:generate sqlc generate

func New(ctx context.Context, db *sql.DB) (*repo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	q, err := queries.Prepare(ctx, db)
	if err != nil {
		return nil, err
	}

	return &repo{q}, nil
}

func Must(ctx context.Context, db *sql.DB) *repo {
	r, err := New(ctx, db)
	if err != nil {
		panic(err)
	}
	return r
}

type repo struct {
	q *queries.Queries
	//db *sql.DB
}

func (r *repo) CreateChat(ctx context.Context, in chatrepos.CreateChatInput) error {
	return r.q.CreateChat(ctx, queries.CreateChatParams{
		ID:        in.ID,
		UserOneID: in.CurrentUserID,
		UserTwoID: in.OtherUserID,
		CreatedAt: time.Now(),
	})
}

type getChatsByUserRow queries.GetChatsByUserRow

func (c getChatsByUserRow) toRepoChat() *chatrepos.Chat {
	return &chatrepos.Chat{
		ID: c.Chat.ID,
		CurrentUser: chatrepos.User{
			ID:        c.User.ID,
			ImageURL:  c.User.ImageUrl.String,
			FirstName: c.User.FirstName,
			LastName:  c.User.LastName.String,
			Username:  c.User.Username,
		},
		OtherUser: chatrepos.User{
			ID:        c.User_2.ID,
			ImageURL:  c.User_2.ImageUrl.String,
			FirstName: c.User_2.FirstName,
			LastName:  c.User_2.LastName.String,
			Username:  c.User_2.Username,
		},
	}
}

func (r *repo) GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*chatrepos.Chat, error) {
	dbChats, err := r.q.GetChatsByUser(ctx, userID)
	if err == sql.ErrNoRows {
		return nil, chatrepos.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var chats []*chatrepos.Chat
	for _, c := range dbChats {
		chats = append(chats, getChatsByUserRow(c).toRepoChat())
	}

	return chats, nil
}

func (r *repo) GetChat(ctx context.Context, chatID uuid.UUID) (*chatrepos.Chat, error) {
	dbC, err := r.q.GetChat(ctx, chatID)
	if err == sql.ErrNoRows {
		return nil, chatrepos.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return getChatsByUserRow(dbC).toRepoChat(), nil
}

func (r *repo) Close() error {
	return r.q.Close()
}

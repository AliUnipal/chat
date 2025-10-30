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
		chats = append(chats, c.ToRepoChat())
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

	return queries.GetChatsByUserRow(dbC).ToRepoChat()., nil
}

func (r *repo) Close() error {
	return r.q.Close()
}

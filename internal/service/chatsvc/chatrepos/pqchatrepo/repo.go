package pqchatrepo

import (
	"context"
	"database/sql"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos/pqchatrepo/queries"
	"github.com/google/uuid"
	"time"
)

//go:generate sqlc generate

func New(ctx context.Context, db *sql.DB) (*repo, error) {
	q, err := queries.Prepare(ctx, db)

	if err != nil {
		return nil, err
	}

	return &repo{
		q:  q,
		db: db,
	}, nil
}

func Must(ctx context.Context, db *sql.DB) *repo {
	r, err := New(ctx, db)
	if err != nil {
		panic(err)
	}
	return r
}

type repo struct {
	q  *queries.Queries
	db *sql.DB
}

func (r *repo) CreateChat(ctx context.Context, in chatrepos.CreateChatInput) error {
	return r.q.CreateChat(ctx, queries.CreateChatParams{
		ID:         in.ID,
		UserOneID:  in.CurrentUserID,
		UserTwoID:  in.OtherUserID,
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	})
}

func dbChatToRepoChat(c queries.GetChatsByUserRow) *chatrepos.Chat {
	var u1ImgUrl sql.NullString
	if c.UserOneImageUrl.String != "" {
		u1ImgUrl = c.UserOneImageUrl
		u1ImgUrl.Valid = true
	}
	var u2ImgUrl sql.NullString
	if c.UserOneImageUrl.String != "" {
		u2ImgUrl = c.UserTwoImageUrl
		u2ImgUrl.Valid = true
	}

	var u1LName sql.NullString
	if c.UserOneLastName.String != "" {
		u1LName = c.UserOneLastName
		u1LName.Valid = true
	}
	var u2LName sql.NullString
	if c.UserOneLastName.String != "" {
		u2LName = c.UserTwoLastName
		u2LName.Valid = true
	}

	return &chatrepos.Chat{
		ID: c.ChatID,
		CurrentUser: chatrepos.User{
			ID:        c.UserOneID,
			ImageURL:  u1ImgUrl.String,
			FirstName: c.UserOneFirstName,
			LastName:  u1LName.String,
			Username:  c.UserOneUsername,
		},
		OtherUser: chatrepos.User{
			ID:        c.UserTwoID,
			ImageURL:  u2ImgUrl.String,
			FirstName: c.UserTwoFirstName,
			LastName:  u2LName.String,
			Username:  c.UserTwoUsername,
		},
	}
}

func (r *repo) GetChatsByUser(ctx context.Context, userID uuid.UUID) ([]*chatrepos.Chat, error) {
	dbChats, err := r.q.GetChatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var chats []*chatrepos.Chat
	for _, c := range dbChats {
		chats = append(chats, dbChatToRepoChat(c))
	}

	return chats, nil
}

func (r *repo) GetChat(ctx context.Context, chatID uuid.UUID) (*chatrepos.Chat, error) {
	dbC, err := r.q.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return dbChatToRepoChat(queries.GetChatsByUserRow(dbC)), nil
}

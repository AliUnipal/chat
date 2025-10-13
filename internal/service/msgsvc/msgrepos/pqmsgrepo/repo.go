package pqmsgrepo

import (
	"context"
	"database/sql"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos/pqmsgrepo/queries"
	"github.com/google/uuid"
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

func (r *repo) CreateMessage(ctx context.Context, in msgrepos.CreateMessageInput) error {
	return r.q.CreateMessage(ctx, queries.CreateMessageParams{
		ID:        in.ID,
		SenderID:  in.SenderID,
		ChatID:    in.ChatID,
		Content:   in.Content,
		CreatedAt: in.Timestamp,
	})
}

func (r *repo) GetMessages(ctx context.Context, chatID uuid.UUID) ([]msgrepos.Message, error) {
	dbMsgs, err := r.q.GetMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}

	var msgs []msgrepos.Message
	for _, m := range dbMsgs {
		msgs = append(msgs, msgrepos.Message{
			ID:       m.ID,
			SenderID: m.SenderID,
			ChatID:   m.ChatID,
			Content:  m.Content,
		})
	}

	return msgs, nil
}

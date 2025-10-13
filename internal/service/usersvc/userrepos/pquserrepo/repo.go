package pquserrepo

import (
	"context"
	"database/sql"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos/pquserrepo/queries"
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

func (r *repo) CreateUser(ctx context.Context, in userrepos.CreateUserInput) error {
	var imgUrl sql.NullString
	if in.ImageURL != "" {
		imgUrl.String = in.ImageURL
		imgUrl.Valid = true
	}
	var lastName sql.NullString
	if in.LastName != "" {
		lastName.String = in.LastName
		lastName.Valid = true
	}

	return r.q.CreateUser(ctx, queries.CreateUserParams{
		ID:         in.ID,
		ImageUrl:   imgUrl,
		FirstName:  in.FirstName,
		LastName:   lastName,
		Username:   in.Username,
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	})
}

func (r *repo) GetUser(ctx context.Context, id uuid.UUID) (userrepos.User, error) {
	u, err := r.q.GetUser(ctx, id)
	if err != nil {
		return userrepos.User{}, err
	}

	var imgUrl string
	if u.ImageUrl.Valid {
		imgUrl = u.ImageUrl.String
	}
	var lastName string
	if u.LastName.Valid {
		lastName = u.LastName.String
	}

	return userrepos.User{
		ID:        u.ID,
		ImageURL:  imgUrl,
		FirstName: u.FirstName,
		LastName:  lastName,
		Username:  u.Username,
	}, nil
}

func (r *repo) Close() error {
	return r.db.Close()
}

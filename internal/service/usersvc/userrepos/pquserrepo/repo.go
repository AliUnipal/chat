package pquserrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos/pquserrepo/queries"
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
	fmt.Println(in.PasswordHash)

	return r.q.CreateUser(ctx, queries.CreateUserParams{
		ID:        in.ID,
		ImageUrl:  imgUrl,
		FirstName: in.FirstName,
		LastName:  lastName,
		Username:  in.Username,
		Password:  in.PasswordHash,
		CreatedAt: time.Now(),
	})
}

func (r *repo) GetUser(ctx context.Context, id uuid.UUID) (userrepos.User, error) {
	u, err := r.q.GetUser(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return userrepos.User{}, userrepos.ErrNotFound
	}
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

func (r *repo) GetUserByUsername(ctx context.Context, username string) (userrepos.User, error) {
	u, err := r.q.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return userrepos.User{}, userrepos.ErrNotFound
	}
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
		ID:           u.ID,
		ImageURL:     imgUrl,
		FirstName:    u.FirstName,
		LastName:     lastName,
		Username:     u.Username,
		PasswordHash: u.Password,
	}, nil
}

func (r *repo) Close() error {
	return r.q.Close()
}

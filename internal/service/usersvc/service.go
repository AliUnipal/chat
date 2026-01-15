package usersvc

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	ImageURL  string
	FirstName string
	LastName  string
	Username  string
	Password  string
}

type AuthenticateOutput struct {
	User  user.User
	Token string
}

type userService interface {
	CreateUser(ctx context.Context, in CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
	Authenticate(ctx context.Context, username, password string) (AuthenticateOutput, error)
}

type userRepository interface {
	CreateUser(ctx context.Context, in userrepos.CreateUserInput) error
	GetUser(ctx context.Context, id uuid.UUID) (userrepos.User, error)
	GetUserByUsername(ctx context.Context, username string) (userrepos.User, error)
}

type hasher interface {
	Hash(s string) ([]byte, error)
	Verify(s string, hash []byte) (bool, error)
}

type iJWTManager interface {
	CreateToken(ctx context.Context, u user.User) (string, error)
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
}

type service struct {
	repo   userRepository
	hasher hasher
	jwt    iJWTManager
}

func NewService(repo userRepository, hasher hasher, jwt iJWTManager) *service {
	return &service{repo, hasher, jwt}
}

var _ userService = (*service)(nil)

func (s *service) CreateUser(ctx context.Context, in CreateUserInput) (uuid.UUID, error) {
	if in.FirstName == "" {
		return uuid.Nil, errors.New("first name is required")
	}
	if in.Username == "" {
		return uuid.Nil, errors.New("username is required")
	}
	if in.Password == "" {
		return uuid.Nil, errors.New("password is required")
	}
	if in.ImageURL == "" {
		return uuid.Nil, errors.New("image url is required")
	}
	u, err := url.ParseRequestURI(in.ImageURL)
	if err != nil || u == nil || u.Scheme == "" || u.Host == "" {
		return uuid.Nil, errors.New("image url is invalid")
	}

	userID := uuid.New()
	hashPass, err := s.hasher.Hash(in.Password)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.repo.CreateUser(ctx, userrepos.CreateUserInput{
		ID:           userID,
		ImageURL:     in.ImageURL,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Username:     in.Username,
		PasswordHash: hashPass,
	}); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:        u.ID,
		ImageURL:  u.ImageURL,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
	}, nil
}

func (s *service) Authenticate(ctx context.Context, username, password string) (AuthenticateOutput, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
	fmt.Printf("%+v\n", u)

	if errors.Is(err, userrepos.ErrNotFound) {
		_, err := s.hasher.Verify(password, nil)
		if err != nil {
			return AuthenticateOutput{}, err
		}

		return AuthenticateOutput{}, errors.New("invalid credentials")
	}
	if err != nil {
		return AuthenticateOutput{}, err
	}

	isValid, err := s.hasher.Verify(password, u.PasswordHash)
	if err != nil {
		return AuthenticateOutput{}, err
	}
	if !isValid {
		return AuthenticateOutput{}, errors.New("invalid credentials")
	}

	usr := user.User{
		ID:        u.ID,
		ImageURL:  u.ImageURL,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
	}

	token, err := s.jwt.CreateToken(ctx, usr)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	return AuthenticateOutput{
		User:  usr,
		Token: token,
	}, nil
}

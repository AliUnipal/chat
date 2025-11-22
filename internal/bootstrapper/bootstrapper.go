package bootstrapper

import (
	"context"
	"database/sql"

	"github.com/AliUnipal/chat/internal/models/chat"
	"github.com/AliUnipal/chat/internal/models/message"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/service/chatsvc"
	"github.com/AliUnipal/chat/internal/service/chatsvc/chatrepos/pqchatrepo"
	"github.com/AliUnipal/chat/internal/service/msgsvc"
	"github.com/AliUnipal/chat/internal/service/msgsvc/msgrepos/pqmsgrepo"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/AliUnipal/chat/internal/service/usersvc/userrepos/pquserrepo"
	"github.com/google/uuid"
)

// New initializes a new bootstrapper with a database connection using the provided database URL.
// Returns a pointer to the bootstrapper or an error if the connection or ping fails.
// The caller is responsible for closing the connection.
func New(dbUrl string) (*bootstrapper, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &bootstrapper{
		dbConn: conn,
	}, nil
}

type bootstrapper struct {
	dbConn *sql.DB
}

func (bs *bootstrapper) Close() error {
	return bs.dbConn.Close()
}

type userService interface {
	CreateUser(ctx context.Context, in usersvc.CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type chatService interface {
	CreateChat(ctx context.Context, currentUserID, otherUserID uuid.UUID) (uuid.UUID, error)
	GetChats(ctx context.Context, userID uuid.UUID) ([]chat.Chat, error)
}

type messageService interface {
	CreateMessage(ctx context.Context, in msgsvc.MessageInput) (uuid.UUID, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]message.Message, error)
}

func (bs *bootstrapper) NewChatService(ctx context.Context) chatService {
	chatRepo := pqchatrepo.Must(ctx, bs.dbConn)
	return chatsvc.NewService(chatRepo)
}

func (bs *bootstrapper) NewUserService(ctx context.Context) userService {
	userSrv := pquserrepo.Must(ctx, bs.dbConn)
	return usersvc.NewService(userSrv)
}

func (bs *bootstrapper) NewMessageService(ctx context.Context) messageService {
	msgRepo := pqmsgrepo.Must(ctx, bs.dbConn)
	return msgsvc.NewService(msgRepo)
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
	_ "github.com/lib/pq"
)

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

type application struct {
	chatSvc    chatService
	userSvc    userService
	messageSvc messageService
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	conn, err := sql.Open("postgres", dbUrl)
	defer func() {
		err := conn.Close()
		fmt.Println(err)
	}()

	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	chatRepo := pqchatrepo.Must(ctx, conn)
	chatSvc := chatsvc.NewService(chatRepo)
	defer func() {
		if err := chatRepo.Close(); err != nil {
			log.Println(err)
		}
	}()

	userRepo := pquserrepo.Must(ctx, conn)
	userSvc := usersvc.NewService(userRepo)
	defer func() {
		if err := userRepo.Close(); err != nil {
			log.Println(err)
		}
	}()

	msgRepo := pqmsgrepo.Must(ctx, conn)
	msgSvc := msgsvc.NewService(msgRepo)
	defer func() {
		if err := msgRepo.Close(); err != nil {
			log.Println(err)
		}
	}()

	app := &application{
		chatSvc:    chatSvc,
		userSvc:    userSvc,
		messageSvc: msgSvc,
	}

	if err := app.handler(ctx); err != nil {
		panic(err)
	}
}

func (app *application) handler(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("arguments must be above one")
	}
	cmd := os.Args[1]
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	if cmd == "" {
		return errors.New("empty command")
	}

	switch cmd {
	case "help":
		fmt.Println("Commands:")
		fmt.Println("- create-user")
		fmt.Println("- get-user")
		fmt.Println("- create-chat")
		fmt.Println("- get-chats-by-user")
		fmt.Println("- send-message")
		fmt.Println("- get-messages")
	case "create-user":
		userInput, err := getCreateUserInput(os.Args[2:])
		if err != nil {
			return err
		}

		userID, err := app.userSvc.CreateUser(ctx, userInput)
		if err != nil {
			return err
		}

		fmt.Printf("User created:%+v\n", user.User{
			ID:        userID,
			ImageURL:  userInput.ImageURL,
			FirstName: userInput.FirstName,
			LastName:  userInput.LastName,
			Username:  userInput.Username,
		})
	case "get-user":
		id, err := getGetUserInput(os.Args[2:])
		if err != nil {
			return err
		}

		u, err := app.userSvc.GetUser(ctx, id)
		if err != nil {
			return err
		}

		fmt.Printf(
			"User Details:\n\nName: %s %s\nImage URL: %s\nUsername: %s\n",
			u.FirstName, u.LastName, u.ImageURL, u.Username,
		)
	case "create-chat":
		currID, otherID, err := getCreateChatInput(os.Args[2:])
		if err != nil {
			return err
		}

		chatID, err := app.chatSvc.CreateChat(ctx, currID, otherID)
		if err != nil {
			return err
		}

		fmt.Println("Created chat ID:", chatID)
	case "get-chats-by-user":
		userID, err := getUserID(os.Args[2:])
		if err != nil {
			return err
		}

		chats, err := app.chatSvc.GetChats(ctx, userID)
		if err != nil {
			return err
		}

		for i, c := range chats {
			fmt.Println("------")
			fmt.Println("Chat No.:", i)
			fmt.Printf("Chat ID: %s\n", c.ID)
			fmt.Printf("User one: %+v\n", c.CurrentUser)
			fmt.Printf("User two: %+v\n", c.OtherUser)
		}
	case "send-message":
		msg, err := getSendMsgInput(os.Args[2:])
		if err != nil {
			return err
		}

		msgID, err := app.messageSvc.CreateMessage(ctx, msg)
		if err != nil {
			return err
		}

		fmt.Println("Created message ID:", msgID)
	case "get-messages":
		chatID, err := getChatIDInput(os.Args[2:])
		if err != nil {
			return err
		}

		msgs, err := app.messageSvc.GetMessages(ctx, chatID)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			fmt.Println("No messages")
		}

		for i, m := range msgs {
			fmt.Println("------")
			fmt.Println("Message No.:", i)
			fmt.Printf("Message ID: %s\n", m.ID)
			fmt.Printf("SenderID: %s\n", m.SenderID)
			fmt.Printf("Content: %s\n", string(m.Content))
			fmt.Printf("Content Type: %v\n", m.ContentType)
			fmt.Printf("Timestamp: %v\n", m.Timestamp)
		}
	default:
		return errors.New("unknown command")
	}

	return nil
}
